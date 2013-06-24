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
        var details = new N.views.MapPopup({
                model: parcel
            }),
            popup = L.popup().setLatLng(parcel.getCoordinates())
                .openOn(this.map);
        $(popup._contentNode).append(details.render().$el);
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
            this.model.get('selectedParcels').fetch({
                data: latlon,
                reset: true
            });
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

    N.views.MapPopup = Backbone.View.extend({
        initialize: function() {
            var view = this;
            view.tmpl = N.app.tmpl['template-parcel-detail'];
        }, 

        render: function() {
            var view = this;
            view.$el.append(view.tmpl({parcel: view.model.toJSON()}));

            view.saveTo = new N.views.CollectionSelect({
                collection: N.app.collections.myCollections,
                parcelId: view.model.id
            }).on('collectionAdd', function(pc) {
                pc.addParcel(view.model);
            }).on('collectionRemove', function(pc) {
                pc.removeParcel(view.model);
            });

            view.$el.append(view.saveTo.$el);
            return view;
        },

        close: function() {
            this.remove();
            this.saveTo.off();
        }
    });

}(Prc));
