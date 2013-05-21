(function(N) {
    N.models.Search = Backbone.Model.extend({
        defaults: {
            searchParcel: null
        }
    });

    N.views.Search = Backbone.View.extend({
        initialize: function() {
            var view = this;
            view.setElement(this.options.searchBar);
            view.$('input').philaddress({
                onSelect: function(parcel) {
                    view.model.set('searchParcel', parcel);    
                }
            });  
        }
    });

}(Prc));
