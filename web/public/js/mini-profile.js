(function() {
    /**
     * @param {jQuery} $post
     */
    MemoApp.MiniProfile = function($post) {
        var $names = $post.find(".mini-profile-name");
        var loadingHtml = "<span class='glyphicon glyphicon-refresh spinning'></span> Loading...";
        $names.each(function() {
            var $name = $(this);
            var $miniProfile = $name.find(".mini-profile");
            if (!$miniProfile) {
                return;
            }
            var profileHash = $name.attr("data-profile-hash");
            $name.hover(function() {
                $miniProfile.html(loadingHtml);
                $miniProfile.show();
                $.ajax({
                    url: MemoApp.URL.ProfileMini + "/" + profileHash,
                    success: function(html) {
                        $miniProfile.html(html);
                    }
                })
            }, function() {
                $miniProfile.hide();
            });
        });
    };
})();