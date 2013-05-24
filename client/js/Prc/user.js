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
            this.model.on('change:loggedIn', this.render);
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
            this.model.on('change:loggedIn', this.render);
            this.setElement(this.options.el);
            this.render();
        }, 

        events: {
            'click button.login': 'login'
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

        login: function() {
            var view = this;
            var l = $.post('/api/v0.1/login/', {
                    username: this.$('input.username').val(),
                    password: this.$('input.password').val()
                },
                function(user) {
                    view.model.set(_.extend({loggedIn: true}, user));
                }
            );

           l.fail(function(e) {
              alert('failed');
           }); 

        }
    });

}(Prc));
