(function (N) {
    N.app = {
        models: {},
        collections: {},
        views: {},
        tmpl: {},
        data: {},

        init: function () {
            new N.TemplateLoader().load(N.app.tmpl);
        }
    };
}(Prc));

