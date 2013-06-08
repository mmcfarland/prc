(function(N) {
    // Is this confusing enough?
    N.models.Collection = Backbone.Model.extend({
        urlRoot: '/api/v0.1/collections/',

        addParcel: function(parcel) {
            var parcelList = this.get('parcels');
            if (!_.contains(parcelList, parcel.id)) {
                parcelList.push(parcel.id);

                $.post({
                    url: this.urlRoot + 'parcels',
                    data: parcel.id
                }).done(function() {

                });
            }
        }
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

    N.views.CollectionSelect = Backbone.View.extend({
        initialize: function() {
            var view = this;
            view.tmpl = N.app.tmpl['template-save-to-collection'];
            view.render();
            view.collection.on('reset', function() {
                view.render();
            });
        }, 

        events: {
            'change select': 'collectionSelected'
        },

        render: function() {
            var view = this,
                $select;
            view.$el.empty().append(view.tmpl());
            _.invoke(view.items, "remove");
            view.items = [];
            $select = view.$('select');
            $select.append(view._makeOption(new N.models.Collection(
                {id: -1, title: 'Select collection'})));
            view.collection.each(function(c) {
                $select.append(view._makeOption(c));
            });
            return this;
        }, 

        collectionSelected: function(e) {
            var cid = parseInt(e.currentTarget.value),
                c = this.collection.findWhere({
                   id: cid 
                });
            this.trigger('collectionChange', c);
        },

        _makeOption: function(c) {
            var item = new N.views.CollectionOption({model: c});
            this.items.push(item);
            return item.$el;
        } 

    });

    N.views.CollectionOption = Backbone.View.extend({
        initialize: function() {
            this.tmpl = N.app.tmpl['template-collection-option'];
            this.render();
        },

        render: function() {
            var $item = $(this.tmpl(this.model.toJSON()));
            this.setElement($item[0]);
        }
    });

}(Prc));
