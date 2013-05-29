(function(N) {
    // Is this confusing enough?
    N.models.Collection = Backbone.Model.extend({
        urlRoot: '/api/v0.1/collection/'
    });

    N.collections.Collections = Backbone.Collection.extend({
        model: N.models.Collection,
        url: '/api/v0.1/collections/'
    });

    N.collections.LocalCollection = N.collections.Collections.extend({
        localStorage: new Backbone.LocalStorage("local-collections")
    });

    N.views.CollectionList = Backbone.View.extend({
        initialize: function() {
           var view = this,m
               list = N.app.tmpl['template-my-collections']();
           view.setElement($(list)[0]);
           view.collection.on('reset', function() {
               view.render();
           });
        },

        render: function() {
            var view = this;
            _.invoke(view.items, "remove");
            view.items = [];
            view.collection.each(function(c) {
                var item = new N.views.CollectionItem({model: c});
                view.items.push(item);
                view.$el.append(item.$el);
            });
            return this;
        }
    });

    N.views.CollectionItem = Backbone.View.extend({
        initialize: function() {
            this.tmpl = N.app.tmpl['template-collection'];
            this.render();
        },

        render: function() {
            var $item = $(this.tmpl(this.model.toJSON()));
            this.setElement($item[0]);
        }
    });
}(Prc));
