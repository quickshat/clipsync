var apiEndpoint = 'http://127.0.0.1';

var app = new Vue({
    el: '#app',
    data: {
        views: ['Settings', 'Connections', 'Logs'],
        selectedView: 'Settings',
        settings: {
            
        },
        interfaces: [],
        logs: []
    },
    mounted: function () {
        apiEndpoint += ':8081';
        this.getSettings();
        this.getInterfaces();
        this.getLogs();
        setInterval(this.getLogs, 500);
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
        }
    }
})