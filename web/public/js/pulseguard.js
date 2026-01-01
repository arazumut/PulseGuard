// Global Scope Functions (Guaranteed availability)
window.showChart = function (serviceId, serviceName) {
    console.log("View clicked for:", serviceName);

    // Show Chart Section
    $('#chart-section').show();
    $('#chart-title').text('Analytics: ' + serviceName);

    // Scroll to chart
    document.getElementById('chart-section').scrollIntoView({ behavior: 'smooth' });

    // Fetch Data
    $.get('/api/v1/services/' + serviceId + '/metrics', function (response) {

        // --- 1. POPULATE STATS ---
        if (response.stats) {
            const up = response.stats.uptime_percentage;
            const avg = response.stats.avg_latency;
            const total = response.stats.total_checks;

            // Uptime
            $('#stat-uptime').text(up.toFixed(2) + '%');
            if (up >= 99.9) $('#stat-uptime').attr('class', 'text-success');
            else if (up >= 95) $('#stat-uptime').attr('class', 'text-warning');
            else $('#stat-uptime').attr('class', 'text-danger');

            // Latency
            $('#stat-latency').text((avg / 1e6).toFixed(2) + ' ms');

            // Total Checks
            $('#stat-checks').text(total);
        }

        // --- 2. DRAW CHART ---
        if (!response.history || response.history.length === 0) {
            $('#latency-chart').html('<p class="text-center text-muted p-4">Not enough data for chart.</p>');
            return;
        }

        // Draw Chart (C3.js)
        const history = response.history.reverse();
        const latencies = ['Latency', ...history.map(d => (d.latency_ns / 1e6).toFixed(2))];

        c3.generate({
            bindto: '#latency-chart',
            data: {
                columns: [latencies],
                type: 'area-spline',
                colors: { 'Latency': '#28a745' } // PulseGuard Green
            },
            axis: {
                x: {
                    type: 'category',
                    categories: history.map(d => new Date(d.checked_at).toLocaleTimeString()),
                    show: false // Hide labels if too crowded
                }
            },
            point: {
                show: false
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
        const safeId = service.id;

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
                    <button class="btn btn-outline-info btn-sm view-btn" 
                            onclick="window.showChart('${safeId}', '${safeName}')">
                        View
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
        if (!data.success) {
            row.find('.status-badge').attr('class', 'status-badge status-critical').text('DOWN');
        }
    }
});
