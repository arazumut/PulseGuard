// --- PULSEGUARD PRO ---

// Global View Functions
window.viewAnalytics = function (serviceId, serviceName) {
    console.log("Loading analytics for:", serviceName);

    const panel = $('#analytics-panel');
    panel.removeClass('hidden');
    $('#chart-service-name').text(serviceName);

    // Reset & Scroll
    $('#stat-uptime').text('Loading...');
    $('#stat-latency').text('--');
    $('#stat-checks').text('--');
    $('#chart-container').html('<div class="text-center p-10 text-gray-500"><i class="fa-solid fa-circle-notch fa-spin text-3xl mb-3"></i><br>Fetching Metrics...</div>');
    panel[0].scrollIntoView({ behavior: 'smooth', block: 'start' });

    // Fetch Data
    $.get(`/api/v1/services/${serviceId}/metrics`)
        .done(function (resp) {
            renderStats(resp.stats);
            renderChart(resp.history);
        })
        .fail(function () {
            toastr.error('Failed to load metrics');
            $('#chart-container').html('<div class="text-center p-10 text-red-400">Error loading data</div>');
        });
};

window.hideAnalytics = function () {
    $('#analytics-panel').addClass('hidden');
};

// --- CRUD Operations ---

window.showAddModal = function () {
    $('#add-modal').removeClass('hidden').addClass('flex');
    $('#input-name').focus();
};

window.hideAddModal = function () {
    $('#add-modal').addClass('hidden').removeClass('flex');
    $('#add-service-form')[0].reset();
};

window.handleCreate = function (e) {
    e.preventDefault();

    // Construct payload
    const payload = {
        name: $('#input-name').val(),
        url: $('#input-url').val(),
        interval: parseInt($('#input-interval').val()),
        slack_enabled: $('#input-slack').is(':checked')
    };

    $.ajax({
        url: '/api/v1/services',
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify(payload),
        success: function () {
            toastr.success('Service created successfully');
            window.hideAddModal();
            setTimeout(() => location.reload(), 1000); // Reload to start monitoring
        },
        error: function (xhr) {
            toastr.error('Failed to create service: ' + (xhr.responseJSON?.error || 'Unknown error'));
        }
    });
};

window.deleteService = function (id) {
    if (!confirm('Are you sure? This will delete all history for this service.')) return;

    $.ajax({
        url: '/api/v1/services/' + id,
        type: 'DELETE',
        success: function () {
            toastr.success('Service deleted');
            $(`#service-${id}`).fadeOut(300, function () { $(this).remove(); });
            // Hide analytics if open for this service
            $('#analytics-panel').addClass('hidden');
        },
        error: function () {
            toastr.error('Failed to delete service');
        }
    });
};

// --- Render Helpers ---

function renderStats(stats) {
    if (!stats) return;
    $('#stat-uptime').text(stats.uptime_percentage.toFixed(2) + '%');
    $('#stat-latency').text((stats.avg_latency / 1e6).toFixed(2) + ' ms');
    $('#stat-checks').text(stats.total_checks.toLocaleString());

    // Color Logic
    const upEl = $('#stat-uptime');
    upEl.removeClass('text-success text-warning text-danger');
    if (stats.uptime_percentage >= 99) upEl.addClass('text-success');
    else if (stats.uptime_percentage >= 95) upEl.addClass('text-warning');
    else upEl.addClass('text-danger');
}

function renderChart(history) {
    if (!history || history.length === 0) {
        $('#chart-container').html('<div class="text-center p-10 text-gray-500"><i class="fa-regular fa-folder-open text-4xl mb-3"></i><br>No historical data available yet.</div>');
        return;
    }

    const dataReverse = history.reverse();
    const latencyValues = ['Latency', ...dataReverse.map(d => (d.latency_ns / 1e6).toFixed(2))];
    const categories = dataReverse.map(d => new Date(d.checked_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' }));

    c3.generate({
        bindto: '#chart-container',
        data: {
            columns: [latencyValues],
            type: 'area-spline',
            colors: { 'Latency': '#3b82f6' }
        },
        axis: {
            x: {
                type: 'category',
                categories: categories,
                show: false
            },
            y: {
                label: 'Latency (ms)',
                padding: { bottom: 0 }
            }
        },
        grid: { y: { show: true } },
        point: { show: false, r: 3 },
        tooltip: { format: { value: function (v) { return v + ' ms'; } } }
    });
}

// --- Main App Initialization ---
$(function () {
    // WebSocket
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;

    let ws;

    function connectWS() {
        ws = new WebSocket(wsUrl);
        ws.onopen = () => toastr.success('Real-time connection established');
        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            updateServiceRow(data);
        };
        ws.onclose = () => {
            console.warn("WS Closed, reconnecting...");
            setTimeout(connectWS, 3000);
        };
    }
    connectWS();

    // Initial Load
    $.get('/api/v1/services', function (response) {
        if (response.data) {
            response.data.forEach(addServiceRow);
        }
    });
});

function addServiceRow(service) {
    if ($(`#service-${service.id}`).length > 0) return;

    const intervalSec = service.interval / 1e9;
    const safeName = service.name.replace(/'/g, "\\'");
    const safeId = service.id;

    let badgeColor = 'bg-gray-600';
    if (service.status === 'HEALTHY') badgeColor = 'bg-green-600';
    if (service.status === 'WARNING') badgeColor = 'bg-yellow-600';
    if (service.status === 'CRITICAL' || service.status === 'DOWN') badgeColor = 'bg-red-600';

    const slackIcon = service.slack_enabled ? '<i class="fa-brands fa-slack text-gray-400 ml-2" title="Slack Enabled"></i>' : '';

    const row = `
        <tr id="service-${service.id}" class="hover:bg-gray-800/50 transition border-b border-gray-800 group">
            <td class="p-4 font-medium text-white">
                ${service.name}
                ${slackIcon}
            </td>
            <td class="p-4 text-gray-400 text-sm font-mono truncate max-w-xs">${service.url}</td>
            <td class="p-4 text-gray-500 text-sm">${intervalSec}s</td>
            <td class="p-4">
                <span class="status-badge px-2 py-1 rounded text-xs font-bold text-white ${badgeColor}">
                    ${service.status}
                </span>
            </td>
            <td class="p-4 text-right font-mono text-primary latency-cell">-</td>
            <td class="p-4 text-right text-sm text-gray-500 time-cell">-</td>
            <td class="p-4 text-center flex justify-center gap-2">
                <button onclick="window.viewAnalytics('${safeId}', '${safeName}')" 
                        class="text-gray-400 hover:text-white bg-gray-700 hover:bg-primary px-3 py-1 rounded-md text-sm transition" title="Analytics">
                    <i class="fa-solid fa-chart-simple"></i>
                </button>
                <button onclick="window.deleteService('${safeId}')" 
                        class="text-gray-400 hover:text-white bg-gray-700 hover:bg-red-600 px-3 py-1 rounded-md text-sm transition" title="Delete">
                    <i class="fa-solid fa-trash"></i>
                </button>
            </td>
        </tr>
    `;
    $('#service-list').append(row);
}

function updateServiceRow(data) {
    const row = $(`#service-${data.service_id}`);
    if (row.length === 0) return;

    row.addClass('bg-gray-700');
    setTimeout(() => row.removeClass('bg-gray-700'), 300);

    row.find('.latency-cell').text((data.latency / 1e6).toFixed(2) + ' ms');
    row.find('.time-cell').text(new Date().toLocaleTimeString());

    const badge = row.find('.status-badge');
    if (!data.success) {
        badge.removeClass('bg-green-600 bg-yellow-600').addClass('bg-red-600').text('DOWN');
    } else {
        if (data.status_code >= 200 && data.status_code < 400 && badge.text().trim() === 'DOWN') {
            badge.removeClass('bg-red-600').addClass('bg-green-600').text('HEALTHY');
        }
    }
}
