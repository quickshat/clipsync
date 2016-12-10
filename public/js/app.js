var apiEndpoint = 'http://127.0.0.1';

var app = new Vue({
    el: '#app',
    data: {
        views: ['Settings', 'Connections', 'Logs'],
        selectedView: 'Settings',
        settings: {

        },
        interfaces: [],
        logs: [],
        connections: {}
    },
    mounted: function () {
        apiEndpoint += ':8081';

        this.getSettings();
        this.getInterfaces();
        this.getLogs();
        this.getConnections();

        setInterval(this.getLogs, 750);
        setInterval(this.getConnections, 750);
    },
    methods: {
        selectTab: function (t) {
            this.selectedView = t;
        },
        getSettings: function () {
            this.$http.get(apiEndpoint + '/local/settings').then((response) => {
                app.settings = response.body;
            }, (response) => { });
        },
        getInterfaces: function () {
            this.$http.get(apiEndpoint + '/local/interfaces').then((response) => {
                app.interfaces = response.body;
            }, (response) => { });
        },
        getLogs: function () {
            this.$http.get(apiEndpoint + '/local/logs').then((response) => {
                app.logs = response.body;
            }, (response) => { });
        },
        getConnections: function () {
            this.$http.get(apiEndpoint + '/local/discoveredDevices').then((response) => {
                app.connections = response.body;
            }, (response) => { });
        },
        saveSettings: function () {
            this.$http.post(apiEndpoint + '/local/settings', this.settings, {
                emulateJSON: true
            }).then((response) => {
                app.getSettings();
            }, (response) => { });
        },
        toggleOnOff: function () {

        },
        setInterface: function (i) {
            this.settings.Interface = i;
            this.saveSettings()
        }
    },
    computed: {
        hasNoConnections: function () {
            for (var prop in this.connections) {
                if (this.connections.hasOwnProperty(prop))
                    return false;
            }

            return JSON.stringify(this.connections) === JSON.stringify({});
        }
    }
})