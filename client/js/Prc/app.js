(function (N) {
    N.app = {
        models: {},
        collections: {},
        views: {},
        tmpl: {},
        data: {},

        init: function () {
            new N.TemplateLoader().load(N.app.tmpl);
            N.app.models.map = new N.models.Map();
            N.app.views.map = new N.views.Map({
                model: N.app.models.map
            });
        }
    };
}(Prc));

