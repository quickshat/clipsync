var app = new Vue({
    el: '#app',
    data: {
        views: ['Settings', 'Connections', 'Logs'],
        selectedView: 'Settings'
    },
    methods: {
        selectTab: function (t) {
            this.selectedView = t;
        }
    }
})