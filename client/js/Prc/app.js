(function (N) {
    N.app = {
        models: {},
        collections: {},
        views: {},
        tmpl: {},
        data: {},

        init: function () {
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
        }, 
        
        setupUser: function() {
            N.app.models.user = new N.models.User(N.bootstrap.user);
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

        }
    };
}(Prc));

