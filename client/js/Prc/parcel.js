(function (N) {
    N.models.Parcel = Backbone.Model.extend({
        urlRoot: '/api/v0.1/parcels/'
    });

    N.collections.Parcels = Backbone.Collection.extend({
        model: N.models.Parcel,
        urlRoot: '/api/v0.1/parcels'
    });

}(Prc));
