(function(N) {
    "use strict"
    N.Prc = _.clone(Backbone.Events);
    N.Prc.models = {};
    N.Prc.collections = {};
    N.Prc.views = {};

    N.Prc.on('error', function(err) {
        console.log("error: " + err);
    });
    
}(this));

$.fn.serializeObject = function()
{
       var o = {};
       var a = this.serializeArray();
       $.each(a, function() {
          if (o[this.name]) {
              if (!o[this.name].push) {
                  o[this.name] = [o[this.name]];
              }
              o[this.name].push(this.value || '');
          } else {
              o[this.name] = this.value || '';
          }
        });
        return o;
};
