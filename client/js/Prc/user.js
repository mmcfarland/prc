(function (N) {
    N.models.User = Backbone.Model.extend({
        defaults: {
            username: null,
            email: null,
            joined: null
        }
    });

    N.views.Login = Backbone.View.extend({
        initialize: function() {
            var view = this;
            this.model.on('change:loggedIn', function() {
                view.render();
            });
            this.setElement(this.options.el);
            this.render();
        },

        render: function() {
            var m, ctx = {action:""};
            if (this.model.get('loggedIn')) {
                ctx.action = (this.model.get('username'));
            } else {
                ctx.action =  "Login";
            }
            m = this.options.tmpl.login(ctx);
            this.$el.empty().append(m);
            return this;
        }
    });


    N.views.LoginForm = Backbone.View.extend({
        initialize: function() {
            var view = this;
            this.model.on('change:loggedIn', function() {
                view.render();
            });
            this.setElement(this.options.el);
            this.render();
        }, 

        events: {
            'click a.login': 'login',
            'click button.logout': 'logout',
            'click input': 'noClose'
        },

        render: function() {
            var form;
            if (this.model.get('loggedIn')) {
                form = this.options.tmpl.logoutForm();
            } else {
                form = this.options.tmpl.loginForm();
            }
            this.$el.empty().append(form);
            return this;
        },

        noClose: function(e) { e.stopPropagation()},

        login: function(e) {
            var view = this;
            var l = $.post('/api/v0.1/login/', {
                    username: this.$('input.username').val(),
                    password: this.$('input.password').val()
                },
                function(user) {
                    view.model.set(_.extend({loggedIn: true}, user));
                }, "json"
            );

           l.fail(function(e) {
              alert('failed');
              console.log(e);
           }); 

           e.preventDefault();
        },

        logout: function(e) {
            var view = this,
                l = $.post('/api/v0.1/logout/')
                    .fail(function(e) { alert('failed')})
                    .done(function() {
                        view.model.set('loggedIn', false);
                    });
        }
    });

}(Prc));
