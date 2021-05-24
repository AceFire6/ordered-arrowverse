(function($, Cookies) {
    'use strict';

    /**
     * key-names for local-storage data
     */
    let localStorageKeys = {
        WATCHED_EPISODES: "watchedEpisodes",
    };

    /**
     * CSS class names used to change visual display of "watched" episodes
     */
    let watchedStateClasses = {
        HIDDEN: 'hidden',
        FAINT: 'faint',
    };

    /**
     * Config values for the current running instance.
     */
    let instanceConfig = {
        watchedEpisodesCssClass: watchedStateClasses.FAINT,
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
        $('.select2-results__option').addClass('dark-mode');
        $('.select2-selection.select2-selection--multiple').addClass('dark-mode');
        $('input.form-control.date-picker').addClass('dark-mode');

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
        $('.select2-results__option').removeClass('dark-mode');
        $('.select2-selection.select2-selection--multiple').removeClass('dark-mode');
        $('input.form-control.date-picker').removeClass('dark-mode');

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

    /**
     * Convert the "data-series" and "data-episode-id" into a sinlge key for
     * LocalStorage.
     * 
     * @param  {Element} element the element which contains the data attributes.
     * @returns A string identifying the episode
     */
    var getLSEpisodeKey = function (element) {
        let series = element.attributes["data-series"].value;
        let episode = element.attributes["data-episode-id"].value;
        let key = `${series}-${episode}`;
        return key;
    }

    /**
     * Updates localStorage with "watched" information of an episode
     * 
     * 
     * @param {ChangeEvent} evt The event 
     */
    var updateWatched = function (evt) {
        let newValue = evt.target.checked;
        let key = getLSEpisodeKey(evt.target);
        let watchedEpisodes = JSON.parse(
            localStorage.getItem(localStorageKeys.WATCHED_EPISODES) || "[]"
        );
        let index = watchedEpisodes.indexOf(key)
        if (newValue === true && index === -1) {
            watchedEpisodes.push(key);
        } else if (newValue === false && index !== -1) {
            watchedEpisodes.splice(index, 1);
        }
        localStorage.setItem(
            localStorageKeys.WATCHED_EPISODES, JSON.stringify(watchedEpisodes)
        );
        setWatchedDisplayState(instanceConfig.watchedEpisodesCssClass);
    };

    /**
     * Hide/Show episodes, according to the "watched" state
     * 
     * @param {String} watchedClass The CSS class to apply for "watched" episodes
     */
    var setWatchedDisplayState = function (watchedClass) {
        let watchedEpisodes = JSON.parse(
            localStorage.getItem(localStorageKeys.WATCHED_EPISODES) || "[]"
        );
        $('.episode').map(function() {
            let key = getLSEpisodeKey(this);
            for (const [_, value] of Object.entries(watchedStateClasses)) {
                $(this).removeClass(value)
            }
            if (watchedEpisodes.includes(key)) {
                $(this).addClass(watchedClass)
            } else {
                $(this).removeClass(watchedClass)
            }
        });
    };

    var registerListeners = function() {
        $('.watchedToggle').change(updateWatched);

        $('#show-watched').click(function() {
            let linkText = "SHOW WATCHED";
            if (instanceConfig.watchedEpisodesCssClass === watchedStateClasses.FAINT) {
                instanceConfig.watchedEpisodesCssClass = watchedStateClasses.HIDDEN;
                linkText = "SHOW WATCHED";
            } else {
                instanceConfig.watchedEpisodesCssClass = watchedStateClasses.FAINT;
                linkText = "HIDE WATCHED";
            }
            setWatchedDisplayState(instanceConfig.watchedEpisodesCssClass);

            // Accessing "firstChild.innerHTML" is brittle. But I deemed this an
            // acceptable trade-off to keep code-churn minimal (unless I missed
            // something). Also, the text-value is decoupled from the HTML
            // template for the same reason. This might lead to subtle
            // display-bugs (only the text-value) if the template is updated.
            this.firstChild.innerHTML = linkText;
        });

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
        setWatchedDisplayState(instanceConfig.watchedEpisodesCssClass);

        var colourSetting = Cookies.get('colour');
        if (colourSetting === undefined) {
            Cookies.set('colour', '1');
        } else if (colourSetting === '0') {
            disableColours();
        }

        var userPrefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
        var darkModeSetting = Cookies.get('dark-mode');

        var darkModeEnabled;
        if (darkModeSetting === undefined) {
            // If no user-defined setting set we defer to media setting
            darkModeEnabled = userPrefersDark;
        } else {
            // otherwise we use user setting
            darkModeEnabled = darkModeSetting === '1';
        }

        if (darkModeEnabled) {
            darkModeEnable();
        } else {
            darkModeDisable();
        }
    });
})(jQuery, Cookies);
