(function (N) {
    N.app = {
        models: {},
        collections: {},
        views: {},
        tmpl: {},
        data: {},

        init: function () {
            N.TemplateLoader().load(N.app.tmpl);
            var m = new N.Map();
        }
    };
}(Prc));

