/**
 * PulseGuard Dashboard Application
 * Architecture: Senior Modular
 */

class ChartManager {
    constructor(elementId) {
        this.elementId = elementId;
        this.chart = null;
        this.currentServiceId = null;
    }

    render(serviceName, dataPoints) {
        // dataPoints = [{ latency: ns, time: 'HH:mm:ss' }]
        const latencies = ['Latency (ms)', ...dataPoints.map(d => (d.latency_ns / 1e6).toFixed(2))];
        // C3.js doesn't handle time axis easily without formatting, let's just use categories or index for simplicity first
        // Or if we want real time axis.

        // Let's keep it simple: Line Chart
        if (this.chart) {
            this.chart.destroy();
        }

        this.chart = c3.generate({
            bindto: '#' + this.elementId,
            data: {
                columns: [latencies],
                type: 'area-spline',
                colors: { 'Latency (ms)': '#28a745' }
            },
            axis: {
                x: {
                    type: 'category',
                    categories: dataPoints.map(d => new Date(d.checked_at).toLocaleTimeString()),
                    tick: { count: 10 }
                },
                y: {
                    label: 'Latency'
                }
            }
        });

        $('#chart-title').text(`Latency History: ${serviceName}`);
        $('#chart-section').show();
    }

    update(result) {
        if (!this.chart || this.currentServiceId !== result.service_id) return;

        // Flow: Add new point
        const ms = (result.latency / 1e6).toFixed(2);
        const time = new Date(result.checked_at).toLocaleTimeString();

        this.chart.flow({
            columns: [
                ['Latency (ms)', ms]
            ],
            // keys: { x: time } // C3 flow with categories is tricky, simplest is reload or just append
            // Implementing proper flow requires x axis management.
            // For MVP, reloading data or just letting user refresh for now. 
            // Actually, C3 flow is great but sensitive.
        });
    }
}

class ServiceManager {
    constructor(tableId, apiEndpoint) {
        this.table = $(tableId);
        this.apiEndpoint = apiEndpoint;
        this.services = new Map();
        this.onViewClick = null; // Callback
    }

    async load() {
        try {
            const response = await $.get(this.apiEndpoint);
            if (response.data) {
                response.data.forEach(s => this.upsert(s));
            }
        } catch (e) {
            toastr.error('Failed to load services');
        }
    }

    upsert(service) {
        if (!this.services.has(service.id)) {
            this.createRow(service);
        } else {
            // Update metadata if needed
        }
        this.services.set(service.id, service);
    }

    createRow(service) {
        const rowId = `service-${service.id}`;
        const statusClass = this.getStatusClass(service.status);
        const intervalSec = service.interval / 1e9;

        const html = `
            <tr id="${rowId}">
                <td><strong>${service.name}</strong></td>
                <td><a href="${service.url}" target="_blank" class="text-muted"><small>${service.url}</small></a></td>
                <td>${intervalSec}s</td>
                <td><span class="status-badge ${statusClass}">${service.status}</span></td>
                <td class="latency-cell">-</td>
                <td class="time-cell">-</td>
                <td>
                    <button class="btn btn-sm btn-outline-primary view-btn" data-id="${service.id}" data-name="${service.name}">
                        <i class="fa fa-line-chart"></i> View
                    </button>
                </td>
            </tr>
        `;
        this.table.append(html);

        // Bind Click
        $(`#${rowId} .view-btn`).click(() => {
            if (this.onViewClick) this.onViewClick(service.id, service.name);
        });
    }

    updateStatus(result) {
        // result = { service_id, latency, ... }
        const row = $(`#service-${result.service_id}`);
        if (row.length === 0) return;

        // Latency
        const ms = (result.latency / 1e6).toFixed(2);
        row.find('.latency-cell').text(ms + ' ms');
        row.find('.time-cell').text(new Date(result.checked_at).toLocaleTimeString());

        // Blink
        row.addClass('row-blink');
        setTimeout(() => row.removeClass('row-blink'), 500);

        // Status Badge Logic (Simple)
        const badge = row.find('.status-badge');
        if (!result.success) {
            badge.attr('class', 'status-badge status-critical').text('DOWN');
        } else if (result.status_code >= 500) {
            badge.attr('class', 'status-badge status-critical').text('CRITICAL');
        }
        // Ideally we receive "new_status" from WS to be accurate
    }

    getStatusClass(status) {
        switch (status) {
            case 'HEALTHY': return 'status-healthy';
            case 'WARNING': return 'status-warning';
            case 'CRITICAL': case 'DOWN': return 'status-critical';
            default: return 'status-unknown';
        }
    }
}

class PulseGuardApp {
    constructor() {
        this.serviceManager = new ServiceManager('#service-list', '/api/v1/services');
        this.chartManager = new ChartManager('latency-chart');
        this.ws = null;
    }

    init() {
        // Wire up interactions
        this.serviceManager.onViewClick = async (id, name) => {
            await this.loadServiceHistory(id, name);
        };

        // Load Initial Data
        this.serviceManager.load();

        // Connect WS
        this.connectWS();
    }

    async loadServiceHistory(id, name) {
        toastr.info(`Loading history for ${name}...`);
        this.chartManager.currentServiceId = id; // Set active context

        try {
            const resp = await $.get(`/api/v1/services/${id}/metrics`);
            // Reverse history to show oldest -> newest (left -> right)
            const history = (resp.history || []).reverse();
            this.chartManager.render(name, history);
            // location.href = "#chart-section"; // Scroll to chart
            document.getElementById('chart-section').scrollIntoView({ behavior: 'smooth' });
        } catch (e) {
            toastr.error('Failed to load history');
        }
    }

    connectWS() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws`;
        this.ws = new WebSocket(wsUrl);

        this.ws.onopen = () => {
            console.log('WS Connected');
            toastr.success('Live Monitor Connected');
        };

        this.ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            this.serviceManager.updateStatus(data);

            // If viewing chart for this service, update it
            if (this.chartManager.currentServiceId === data.service_id) {
                // To keep chart simple, maybe just re-fetch in MVP or use flow()
                // this.chartManager.update(data);
            }
        };

        this.ws.onclose = () => {
            toastr.warning('Connection lost. Reconnecting...');
            setTimeout(() => this.connectWS(), 5000);
        };
    }
}

// Bootstrapping
$(function () {
    window.app = new PulseGuardApp();
    window.app.init();
});
