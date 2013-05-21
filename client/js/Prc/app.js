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
        }
    };
}(Prc));

