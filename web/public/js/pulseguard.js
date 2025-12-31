// Global Scope Functions (Guaranteed availability)
window.showChart = function (serviceId, serviceName) {
    console.log("View clicked for:", serviceName);

    // Show Chart Section
    $('#chart-section').show();
    $('#chart-title').text('History: ' + serviceName);

    // Scroll to chart
    document.getElementById('chart-section').scrollIntoView({ behavior: 'smooth' });

    // Fetch Data
    $.get('/api/v1/services/' + serviceId + '/metrics', function (response) {
        if (!response.history) return;

        // Draw Chart (C3.js)
        const history = response.history.reverse();
        const latencies = ['Latency', ...history.map(d => (d.latency_ns / 1e6).toFixed(2))];

        c3.generate({
            bindto: '#latency-chart',
            data: {
                columns: [latencies],
                type: 'area-spline',
                colors: { 'Latency': '#28a745' }
            },
            axis: {
                x: {
                    type: 'category',
                    categories: history.map(d => new Date(d.checked_at).toLocaleTimeString()),
                    show: false // Hide labels if too crowded
                }
            }
        });
    }).fail(function () {
        toastr.error('Could not load history');
    });
};

// Main App Logic
$(function () {
    const table = $('#service-list');

    // 1. WebSocket Connetion
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const ws = new WebSocket(`${protocol}//${window.location.host}/ws`);

    ws.onmessage = function (event) {
        const data = JSON.parse(event.data);
        updateRow(data);
    };

    // 2. Load Initial Services
    $.get('/api/v1/services', function (response) {
        if (response.data) {
            response.data.forEach(s => createRow(s));
        }
    });

    function createRow(service) {
        if ($('#service-' + service.id).length > 0) return;

        const intervalSec = service.interval / 1e9;
        const safeName = service.name.replace(/'/g, "\\'"); // Escape quotes

        // Status Color
        let badgeClass = 'status-unknown';
        if (service.status === 'HEALTHY') badgeClass = 'status-healthy';
        if (service.status === 'WARNING') badgeClass = 'status-warning';
        if (service.status === 'CRITICAL' || service.status === 'DOWN') badgeClass = 'status-critical';

        const html = `
            <tr id="service-${service.id}">
                <td><strong>${service.name}</strong></td>
                <td><a href="${service.url}" target="_blank">${service.url}</a></td>
                <td>${intervalSec}s</td>
                <td><span class="status-badge ${badgeClass}">${service.status}</span></td>
                <td class="latency">-</td>
                <td class="last-check">-</td>
                <td>
                    <button class="btn btn-outline-info btn-sm" 
                            onclick="window.showChart('${service.id}', '${safeName}')">
                        View History
                    </button>
                </td>
            </tr>
        `;
        table.append(html);
    }

    function updateRow(data) {
        const row = $('#service-' + data.service_id);
        if (row.length === 0) return;

        const ms = (data.latency / 1e6).toFixed(2);
        row.find('.latency').text(ms + ' ms');
        row.find('.last-check').text(new Date(data.checked_at).toLocaleTimeString());

        // Blink
        row.addClass('row-blink');
        setTimeout(() => row.removeClass('row-blink'), 500);

        // Update Badge if needed (Optional: handled by reload usually)
    }
});
