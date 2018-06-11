(function () {
    /**
     * @param {jQuery} $form
     */
    MemoApp.Form.Settings = function ($form) {
        var $saved = $form.find("#saved");
        $form.submit(function (e) {
            e.preventDefault();
            var defaultTipRaw = $form.find("[name=default-tip]").val();
            var defaultTip = parseInt(defaultTipRaw);
            var integrations = $form.find("[name=integrations]:checked").val();
            var theme = $form.find("[name=theme]:checked").val();

            if (defaultTipRaw.length > 0) {
                if (isNaN(defaultTip)) {
                    MemoApp.AddAlert("Must enter a numeric default tip.");
                    return;
                }
                if (defaultTip < 0) {
                    MemoApp.AddAlert("Cannot have a negative tip value.");
                    return;
                }
                if (defaultTip !== 0 && defaultTip < 546) {
                    MemoApp.AddAlert("Default tip must be above dust limit of 546 satoshis.")
                    return;
                }
            }
            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + MemoApp.URL.ProfileSettingsSubmit,
                data: {
                    defaultTip: defaultTip,
                    integrations: integrations,
                    theme: theme
                },
                success: function () {
                    $saved.removeClass("hidden");
                    setTimeout(function() {
                        window.location.reload();
                    }, 500);
                },
                /**
                 * @param {XMLHttpRequest} xhr
                 */
                error: function (xhr) {
                    var errorMessage =
                        "Error saving settings:\nCode: " + xhr.responseText + "\n" +
                        "If this problem persists, try refreshing the page.";
                    MemoApp.AddAlert(errorMessage);
                }
            });
        });
    };
})();
