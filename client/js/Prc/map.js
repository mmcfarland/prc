(function (N) {

    N.Map = function(options) {
        var self = this,
            osm = L.tileLayer.provider('MapQuestOpen.OSM');
        options = _.extend({
            el: 'map'
        }, options);

        self.map = new L.map(options.el, {
            center: [39.95, -75.15],
            zoom: 15
        });
      
        // Highlight selected parcels layer
        self.selectedLayer = L.layerGroup().addTo(self.map); 
//      self.map.fitBounds([20,30]);

        self.map.on("click", function(e) {
            self.trigger("click", {latlon: {
                lat: e.latlng.lat, lon: e.latlng.lng
            }});
        });

        self.map.addLayer(osm);
        _.extend(self, Backbone.Events)
    };
    N.Map.prototype.popup = function(parcel) {
        var context = {parcel: parcel.toJSON()},
            content = N.app.tmpl["template-parcel-detail"](context);

        return L.popup().setLatLng(parcel.getCoordinates())
                .setContent(content)
                .openOn(this.map);
    };

    N.Map.prototype.selectSingleGeom = function(geom) {
        L.geoJson(geom).addTo(this.selectedLayer.clearLayers());
    };

    N.models.Map = Backbone.Model.extend({
        defaults: {
            selectedParcels: new N.collections.Parcels(),
            highlightParcel: null
        }, initialize: function() { var self = this; }
    });

    /* VIEW */
    N.views.Map = Backbone.View.extend({

        initialize: function() {
            var self = this;
            self._map = new N.Map();
            self._map.on('click', function(e) {
                self.makeParcelInfoRequest(e.latlon);
            });

            self.model.get('selectedParcels').on('reset', function() {
                // `this` is the selectedParcels collection
                self._handleSelectedParcels(this);    
            });

            self.options.search.on('change:searchParcel', function(s, pi) {
                self.showParcel(pi);
            });
        },

        makeParcelInfoRequest: function(latlon) {
            this.model.get('selectedParcels').fetch({data: latlon});
        }, 

        showParcel: function(parcelInfo) {
            var p = new N.models.Parcel({id: parcelInfo.parcelId}),
                view = this;

            p.fetch().done(function(d) {
                view.model.get('selectedParcels').reset([p]);
             });
        },

        _handleSelectedParcels: function(parcels) {
            var parcel = parcels.at(0);
            this.popup = this._map.popup(parcel);
            this._map.selectSingleGeom(parcel.getGeom());
        }

            
    });

    N.views.Popup = Backbone.View.extend({
        initialize: function() {
            var view = this;
            view.tmpl = N.app.tmpl['template-popup'];
        }, 

        render: function() {

            var collection = new N.views.CollectionSelect({
                collection: N.app.collections.myCollections
            });
        }
    });

}(Prc));
