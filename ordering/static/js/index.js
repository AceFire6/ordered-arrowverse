(function($) {
    'use strict';

    var openWiki = function() {
        var url = $(this).data('href');
        window.open(url, '_blank');
    };

    var registerListeners = function() {
        $('.episode').click(openWiki);
    };

    $(document).ready(registerListeners);
})(jQuery);
