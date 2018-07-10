(function () {
    var $modalWrapper,
        $siteCover,
        $siteWrapper,
        $innerWrapper,
        $document,
        $window,
        windowHeight,
        scrollTop;
    $(function () {
        $modalWrapper = $("#site-modal-wrapper");
        $siteCover = $("#site-wrapper-cover");
        $siteWrapper = $("#site-wrapper");
        $innerWrapper = $("#inner-site-wrapper");
        $document = $(document);
        $window = $(window);
    });
    /**
     * @param {string} html
     */
    MemoApp.Modal = function (html) {
        html =
            "<div class='container'>" +
            html +
            "</div>";
        $modalWrapper.html(html).show();
        windowHeight = $window.height();
        console.log(windowHeight);
        scrollTop = $document.scrollTop();
        $siteWrapper.css({height: windowHeight});
        $innerWrapper.css({top: -scrollTop});
        $siteWrapper.addClass("active");
        $siteCover.addClass("active");
        $document.scrollTop(0);
        $modalWrapper.find(".container").click(function (e) {
            e.stopPropagation();
        });
        $modalWrapper.click(function () {
            $siteWrapper.removeClass("active");
            $innerWrapper.css({top: 0});
            $document.scrollTop(scrollTop);
            $siteCover.removeClass("active");
            $modalWrapper.hide();
        });
    };
})();
