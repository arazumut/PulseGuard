$(function () {
    const serviceList = $('#service-list');
    const servicesMap = new Map(); // ID -> Service Object

    // 1. Initial Load
    function loadServices() {
        $.get('/api/v1/services', function (response) {
            if (response.data) {
                response.data.forEach(service => {
                    addOrUpdateServiceRow(service);
                });
            }
        });
    }

    function addOrUpdateServiceRow(service) {
        let row = $(`#service-${service.id}`);
        const statusClass = getStatusClass(service.status);

        const html = `
            <td><strong>${service.name}</strong></td>
            <td><a href="${service.url}" target="_blank">${service.url}</a></td>
            <td>${service.interval / 1000000000}s</td>
            <td><span class="status-badge ${statusClass}">${service.status}</span></td>
            <td class="latency-cell">-</td>
            <td class="time-cell">-</td>
        `;

        if (row.length === 0) {
            // Create new
            row = $(`<tr id="service-${service.id}">${html}</tr>`);
            serviceList.append(row);
        } else {
            // Update existing (meta data update logic if needed)
            // Usually we only update status via WS, but full refresh handles this too
        }
        servicesMap.set(service.id, service);
    }

    function getStatusClass(status) {
        switch (status) {
            case 'HEALTHY': return 'status-healthy';
            case 'WARNING': return 'status-warning';
            case 'CRITICAL': case 'DOWN': return 'status-critical';
            default: return 'status-unknown';
        }
    }

    function formatLatency(ns) {
        const ms = ns / 1e6;
        return ms.toFixed(2) + ' ms';
    }

    // 2. WebSocket Connection
    function connectWS() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws`;
        const socket = new WebSocket(wsUrl);

        socket.onopen = function () {
            console.log('PulseGuard WS Connected');
            toastr.success('Connected to Live Monitor');
        };

        socket.onmessage = function (event) {
            const data = JSON.parse(event.data);
            // data is CheckResult {service_id, status_code, latency, success...}
            updateServiceStatus(data);
        };

        socket.onclose = function () {
            console.log('WS Disconnected. Reconnecting in 5s...');
            toastr.error('Connection Lost. Reconnecting...');
            setTimeout(connectWS, 5000);
        };
    }

    function updateServiceStatus(result) {
        const row = $(`#service-${result.service_id}`);
        if (row.length > 0) {
            // Update Latency
            const latencyText = formatLatency(result.latency);
            row.find('.latency-cell').text(latencyText);

            // Update Time
            const timeStr = new Date(result.checked_at).toLocaleTimeString();
            row.find('.time-cell').text(timeStr);

            // Blink effect
            row.addClass('row-blink');
            setTimeout(() => row.removeClass('row-blink'), 1000);

            // Note: The 'status' field (HEALTHY/WARNING) comes from Analyzer update,
            // but CheckResult only has raw data (latency, success).
            // To be perfectly sync, we might need the Analyzer to broadcast the NEW STATUS too.
            // For now, let's guess status or wait for full refresh.
            // BETTER: Let the WS send "ServiceUpdate" event not just "CheckResult".

            // Temporary Logic for Status Badge based on Success
            const badge = row.find('.status-badge');
            if (!result.success) {
                badge.removeClass().addClass('status-badge status-critical').text('DOWN');
            } else if (result.status_code >= 500) {
                badge.removeClass().addClass('status-badge status-critical').text('CRITICAL');
            } else {
                // We don't know if it's WARNING or HEALTHY without the threshold logic here.
                // Ideally Backend sends the new derived status.
                // For now let's keep it simple: if success -> Healthy (or keep previous)
            }
        }
    }

    // Initialize
    loadServices();
    connectWS();
});
