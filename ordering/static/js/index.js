(function($) {
    'use strict';

    var openWiki = function() {
        var url = $(this).data('href');
        window.open(url, '_blank');
    };

    var disableColours = function() {
        $('.episode, thead').addClass('no-color');
        $('#episode-list').addClass('table-striped table-hover');
        $('#no-color').find('.text').text('ENABLE COLOR');
        Cookies.set('colour', '0');
    };

    var enableColours = function() {
        $('.episode, thead').removeClass('no-color');
        $('#episode-list').removeClass('table-striped table-hover');
        $('#no-color').find('.text').text('DISABLE COLOR');
        Cookies.set('colour', '1');
    };

    var registerListeners = function() {
        $('.episode').click(openWiki);

        $('#no-color').click(function() {
            if (Cookies.get('colour') === '1') {
                disableColours();
            } else {
                enableColours();
            }
        })
    };

    $(document).ready(function() {
        $('#show-filter-select').select2({
          placeholder: 'Select shows to exclude...',
          allowClear: true,
          closeOnSelect: false,
          width: "100%",
        });
        registerListeners();

        var colourSetting = Cookies.get('colour');
        if (colourSetting === undefined) {
            Cookies.set('colour', '1');
        } else if (colourSetting === '0') {
            disableColours();
        }
    });
})(jQuery);
