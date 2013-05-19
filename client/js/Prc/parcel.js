(function (N) {
    N.models.Parcel = Backbone.Model.extend({
        urlRoot: '/api/v0.1/parcels/',
        getCoordinates: function() {
            if (!this._coords) {
               this._coords = $.parseJSON(this.get('Pos')).coordinates.reverse();
            }
            return this._coords;
        },

        getGeom: function() {
            if (!this._geom) {
                this._geom = $.parseJSON(this.get('Geom'));
            }
            return this._geom;
        }

    });

    N.collections.Parcels = Backbone.Collection.extend({
        model: N.models.Parcel,
        url: '/api/v0.1/parcels/'
    });

}(Prc));
