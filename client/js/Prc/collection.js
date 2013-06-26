(function(N) {
    // Is this confusing enough?
    N.models.Collection = Backbone.Model.extend({
        urlRoot: '/api/v0.1/collections/',

        _addOrRemoveParcel: function(action, parcel, parcelList) {
            var model = this;
            $.ajax({
                type: action,
                url: model.url() + '/parcels/' + parcel.id,
            }).done(function() {
                console.log(arguments);
                // Force the change event by adding a new array, not
                // changing the existing
                model.set('parcelIds', parcelList);
            });
        },

        addParcel: function(parcel) {
            var parcelList = _.clone(this.get('parcelIds'));
            if (!_.contains(parcelList, parcel.id)) {
                parcelList.push(parcel.id);
                this._addOrRemoveParcel('PUT', parcel, parcelList);
            }
        },

        removeParcel: function(parcel) {
            var parcelList = _.clone(this.get('parcelIds'));
            if (_.contains(parcelList, parcel.id)) {
                parcelList.splice(parcelList.indexOf(parcel.id), 1);
                this._addOrRemoveParcel('DELETE', parcel, parcelList);
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
           view.collection.on('reset add', function() {
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
        tagName: "li",
        className: "collection",

        initialize: function() {
            var view = this;
            view.model.on('change:parcelIds', function() {
                view.render();
            });
            view.tmpl = N.app.tmpl['template-collection'];
            view.render();
        },

        render: function() {
            var $item = $(this.tmpl(this.model.toJSON()));
            this.$el.empty().append($item);
        }
    });

    N.views.CollectionSelect = Backbone.View.extend({
        initialize: function() {
            var view = this;
            view.tmpl = N.app.tmpl['template-save-to-collection'];
            view.render();
            view.collection.on('reset sync', function() {
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
            
            $select = view.$('select');
            $select.append(view._makeOption(new N.models.Collection(
                {id: -1, title: 'New Collection...'})));
            $select.append(view._makeOption(new N.models.Collection(
                {id: -2, title: 'Local Scratchpad'})));
            
            view.collection.each(function(c) {
                $select.append(view._makeOption(c));
            });
            $select.chosen();
            return this;
        }, 

        collectionSelected: function(e, item) {
            // TODO: Make a 'new item' value
            if (item.selected === "-1") {
                
                var n = new N.views.AddCollection({
                    parcelList: [this.options.parcelId]
                }).show();
                return;
            } else if (item.selected === "-2" {
                // Localstorage
            }
            
            var c = this.collection.get(parseInt(item.selected || item.deselected));
            if (item.selected) {
                this.trigger('collectionAdd', c);
            } else {
                this.trigger('collectionRemove', c);
            } 
        },

        _makeOption: function(c) {
            var item = new N.views.CollectionOption({
                model: c, parcelId: this.options.parcelId
            });
            return item.$el;
        } 

    });

    N.views.CollectionOption = Backbone.View.extend({
        initialize: function() {
            this.tmpl = N.app.tmpl['template-collection-option'];
            this.render();
        },

        render: function() {
            var containsParcel = _.contains(this.model.get('parcelIds'), this.options.parcelId),
                ctx = _.extend(this.model.toJSON(),{
                    containsParcel: containsParcel
                }),
                $item = $(this.tmpl(ctx));
            this.setElement($item[0]);
        }
    });

    N.views.AddCollection = Backbone.View.extend({
        initialize: function() {
            this.tmpl = N.app.tmpl['template-new-collection'];
            this.render();
        }, 

        events: {
            'click #add-collection': 'addCollection', 
            'click a.cancel': 'close'
        },

        render: function() {
            this.setElement(this.tmpl());
            $('body').append(this.$el);
            return this;
        }, 

        addCollection: function() {
            var newColl = this.$('form').serializeObject(); 
            if (this.options.parcelList) {
                newColl.parcelIds = this.options.parcelList;
            }
            N.app.collections.myCollections.create(newColl);
            this.close();    
        }, 

        show: function() {
            this.$el.foundation('reveal', 'open');
            this.$('input:first').focus();
        },

        close: function() {
            this.$el.foundation('reveal', 'close');
            this.remove();
        }

    });

}(Prc));
