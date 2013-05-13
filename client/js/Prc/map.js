(function (N) {

    N.Map = function(options) {
        var self = this,
            osm = L.tileLayer.provider('MapQuestOpen.OSM');
        options = _.extend({
            el: 'map'
        }, options);

        self.map = new L.map(options.el, {
            center: [39.95, -75.15],
            zoom: 12
        });
       
        self.map.addLayer(osm);
    };

    N.models.Map = Backbone.Model.extend({
    });

    N.views.Map = Backbone.View.extend({

        initialize: function() {
            this.map = new N.Map();
            this.map.on('click', function(e) {
                makeParcelInfoRequest(e.latlng);
            });
        },

        makeParcelInfoRequest: function() {}
    });


}(Prc));
