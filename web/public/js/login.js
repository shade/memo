(function () {
    /**
     * @param {jQuery} $ele
     */
    MemoApp.Form.Login = function ($ele) {
        $ele.submit(function (e) {
            e.preventDefault();
            var username = $ele.find("[name=username]").val();
            var password = $ele.find("[name=password]").val();

            if (username.length === 0) {
                MemoApp.AddAlert("Must enter a username.");
                return;
            }

            if (password.length === 0) {
                MemoApp.AddAlert("Must enter a password.");
                return;
            }

            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.LoginSubmit,
                data: {
                    username: username,
                    password: password
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
                        case 401:
                            MemoApp.AddAlert("Invalid username or password. Please try again.");
                            return
                        case 500:
                            MemoApp.AddAlert("Server side issue. Please try again.");
                            return
                    }
                    var errorMessage =
                        "Error logging in:\n" + xhr.responseText + "\n" +
                        "If this problem persists, try refreshing the page.";
                    MemoApp.AddAlert(errorMessage);
                }
            });
        });
    };
})();
