(function(N) {
    "use strict"
    N.Prc = _.clone(Backbone.Events);
    N.Prc.models = {};
    N.Prc.collections = {};
    N.Prc.Views = {};

    N.Prc.on('error', function(err) {
        console.log("error: " + err);
    });
    
}(this));
