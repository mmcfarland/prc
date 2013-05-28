(function (N) {
    N.app = {
        models: {},
        collections: {},
        views: {},
        tmpl: {},
        data: {},

        init: function () {
            // Foundation must be initialized for all that ui to work
            $(document).foundation();
            new N.TemplateLoader().load(N.app.tmpl);
            N.app.models.search = new N.models.Search();
            N.app.views.search = new N.views.Search({
                model: N.app.models.search,
                searchBar: document.getElementById("search-container")
            });

            N.app.models.map = new N.models.Map();
            N.app.views.map = new N.views.Map({
                model: N.app.models.map,
                search: N.app.models.search
            });

            this.setupUser();
            this.setupParcelCollections();
        }, 
        
        setupUser: function() {
            var u = N.bootstrap.user && N.bootstrap.user.username ? _.extend(N.bootstrap.user, {loggedIn: true}) : {loggedIn: false};
            N.app.models.user = new N.models.User(u);
            N.app.views.login = new N.views.Login({
                model: N.app.models.user,
                el: $('#login-container')[0],
                tmpl: {
                    login: N.app.tmpl["template-login"]
                }
            });

            N.app.views.loginForm = new N.views.LoginForm({
                model: N.app.models.user,
                tmpl: {
                    loginForm: N.app.tmpl["template-login-form"],
                    logoutForm: N.app.tmpl["template-logout-form"]
                },
                el: $('#login-dropdown')[0]
            });

        },

        setupParcelCollections: function() {
            // Parcel collections come in several varieties:
            //  * My Collections - remotely stored specific to user
            //  * Public/Gallery - remotely stored global
            //  * My Computer    - locally stored specific
            var $container = $("#collection-town");

            N.app.collections.myCollections = new N.collections.Collections();
            N.app.collections.localCollection = new N.collections.LocalCollection();

            N.app.collections.myCollections.reset(N.bootstrap.collections);
            N.app.collections.localCollection.fetch();

            N.views.myCollectionList = new N.views.CollectionList({
                collection: N.app.collections.myCollections
            });

            N.views.myCollectionList.render().$el.appendTo($container);
        }
    };
}(Prc));

