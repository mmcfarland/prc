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
       
        self.map.on("click", function(e) {
            self.trigger("click", {latlon: {
                lat: e.latlng.lat, lon: e.latlng.lng
            }});
        });

        self.map.addLayer(osm);
        _.extend(self, Backbone.Events)
    };
    N.Map.prototype.popup = function(parcel) {
        return L.popup().setLatLng($.parseJSON(parcel.get('Pos')).coordinates.reverse())
                .setContent(parcel.get('Owner1')).openOn(this.map);
    };

    N.models.Map = Backbone.Model.extend({
        defaults: {
            selectedParcels: new N.collections.Parcels(),
            highlightParcel: null
        }, 

        initialize: function() {
            var self = this;
            
        }
    });

    N.views.Map = Backbone.View.extend({

        initialize: function() {
            var self = this;
            self.map = new N.Map();
            self.map.on('click', function(e) {
                self.makeParcelInfoRequest(e.latlon);
            });
            self.model.get('selectedParcels').on('reset', function() {
                self.popup = self.map.popup(this.at(0));
                
            });
        },

        makeParcelInfoRequest: function(latlon) {
            this.model.get('selectedParcels').fetch({data: latlon});
        }
    });


}(Prc));
