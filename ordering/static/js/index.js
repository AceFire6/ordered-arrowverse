(function($, Cookies) {
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

    var darkModeEnable = function() {
        $('nav').removeClass('bg-light')
            .removeClass('navbar-light')
            .addClass('bg-dark')
            .addClass('navbar-dark');
        $('body').addClass('dark-mode');
        $('a').addClass('dark-mode');
        $('.modal-content').addClass('dark-mode');
        $('.btn-close').addClass('dark-mode');

        Cookies.set('dark-mode', '1');
    };

    var darkModeDisable = function() {
        $('nav').addClass('bg-light')
            .addClass('navbar-light')
            .removeClass('bg-dark')
            .removeClass('navbar-dark');
        $('body').removeClass('dark-mode');
        $('a').removeClass('dark-mode');
        $('.modal-content').removeClass('dark-mode');
        $('.btn-close').removeClass('dark-mode');

        Cookies.set('dark-mode', '0');
    };

    var toggleDarkMode = function() {
        var darkModeSetting = Cookies.get('dark-mode');
        if (darkModeSetting === '1') {
            darkModeDisable();
        } else {
            darkModeEnable();
        }
    };

    var registerListeners = function() {
        $('.episode').click(openWiki);

        $('#no-color').click(function() {
            if (Cookies.get('colour') === '1') {
                disableColours();
            } else {
                enableColours();
            }
        });

        $('.date-picker').daterangepicker({
            autoUpdateInput: false,
            showDropdowns: true,
            minDate: '2012-10-10',
            singleDatePicker: true,
            locale: {
                format: 'YYYY-MM-DD',
                cancelLabel: 'Clear'
            }
        }).on('apply.daterangepicker', function(ev, picker) {
            $(this).val(picker.startDate.format('YYYY-MM-DD'));
        }).on('cancel.daterangepicker', function (event, picker) {
            $(event.currentTarget).val('');
        });

        $('#dark-mode').click(toggleDarkMode);
    };

    $(document).ready(function() {
        $('#show-filter-select').select2({
          placeholder: 'Select shows to exclude...',
          allowClear: true,
          closeOnSelect: false,
          width: '100%',
        });
        registerListeners();

        var colourSetting = Cookies.get('colour');
        if (colourSetting === undefined) {
            Cookies.set('colour', '1');
        } else if (colourSetting === '0') {
            disableColours();
        }

        var darkModeSetting = Cookies.get('dark-mode');
        if (darkModeSetting === '1') {
            darkModeEnable();
        } else {
            darkModeDisable();
        }
    });
})(jQuery, Cookies);
