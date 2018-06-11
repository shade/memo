(function () {
    /**
     * @param {jQuery} $form
     * @param {jQuery} $privateKeyField
     */
    MemoApp.Form.Signup = function ($form, $privateKeyField) {
        var $radio = $form.find("[name=key-type]");
        $radio.change(function () {
            console.log($radio.val());
            if ($radio.filter(':checked').val() === "import") {
                $privateKeyField.show();
            } else {
                $privateKeyField.hide();
            }
        });

        // Warn if password is < 13 characters.
        var $password = $form.find("[name=password]");
        var passwordWarnOnKey = false
        var checkPassLength = function() {
            if($password.val().length < 13) {
                $('#password-warning').removeClass('hidden');
            } else {
                $('#password-warning').addClass('hidden');
            }
        }
        $password.on("focusout", function() {
            passwordWarnOnKey = true
            checkPassLength()
        });
        $password.on("keypress", function() {
            if(passwordWarnOnKey) {
                checkPassLength()
            }
        });

        $form.submit(function (e) {
            e.preventDefault();
            var username = $form.find("[name=username]").val();
            var password = $form.find("[name=password]").val();
            var retypePassword = $form.find("[name=retype-password]").val();
            if (!$form.find("[name=accept]").is(':checked')) {
                MemoApp.AddAlert("Please accept the disclaimer");
                return;
            }

            if (username.length === 0) {
                MemoApp.AddAlert("Must enter a username.");
                return;
            }

            if (password.length === 0) {
                MemoApp.AddAlert("Must enter a password.");
                return;
            }

            if (retypePassword.length === 0) {
                MemoApp.AddAlert("Must retype password.");
                return;
            }

            if (retypePassword !== password) {
                MemoApp.AddAlert("Password don't match.");
                return;
            }

            var privateKey;
            if ($radio.filter(':checked').val() === "import") {
                privateKey = $form.find("[name=private-key]").val();
                if (privateKey.length === 0) {
                    MemoApp.AddAlert("Must enter a private key for import.");
                    return;
                }
            }

            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.SignupSubmit,
                data: {
                    username: username,
                    password: password,
                    wif: privateKey
                },
                success: function () {
                    MemoApp.SetPassword(password);
                    window.location = MemoApp.GetBaseUrl() + MemoApp.URL.Index
                },
                /**
                 * @param {XMLHttpRequest} xhr
                 */
                error: function (xhr) {
                    switch(xhr.status) {
                        case 422:
                            MemoApp.AddAlert("Could not parse the WIF. Please check the WIF and try again.");
                            return;
                        case 403:
                            MemoApp.AddAlert("Username is not available. Please try a different username.");
                            return;
                        case 401:
                        case 500:
                            MemoApp.AddAlert(
                                "Server side issue while creating the account. " +
                                "Please try again. " +
                                "If the issue persists please try refreshing the page.");
                            return;
                    }
                    MemoApp.AddAlert("Oops! There was an issue creating your account (code: " + xhr.status + "). " +
                        "Try again? " +
                        "If the issue persists please try refreshing the page.");
                }
            });
        });
    };
})();
