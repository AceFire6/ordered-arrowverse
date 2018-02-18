(function($) {
    'use strict';

    var openWiki = function() {
        var url = $(this).data('href');
        window.open(url, '_blank');
    };

    var addFilters = function() {
        var toFilterList = [];

        var selectedOptions = $('#show-filter-select').select2('data');
        $.each(selectedOptions, function() {
            toFilterList.push(this.id);
        });

        var url = '';

        if (toFilterList.length > 0) {
            url = '/hide/' + toFilterList.join('+');
        }

        url += '/';

        if (document.baseURI.match('/newest_first$')) {
            url += 'newest_first';
        }
        window.location = url;
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
        $('#filter-button').click(addFilters);

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
