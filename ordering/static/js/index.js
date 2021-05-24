(function($) {
    'use strict';

    /**
     * key-names for local-storage data
     */
    let localStorageKeys = {
        WATCHED_EPISODES: "watchedEpisodes",
        CONFIG: "config",
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

    /**
     * Load config-values from localStorage into the "instanceConfig"
     */
    var loadConfig = function () {
        let config = JSON.parse(
            localStorage.getItem(localStorageKeys.CONFIG) || '{}'
        );
        if (config.hideWatched) {
            instanceConfig.watchedEpisodesCssClass = watchedStateClasses.HIDDEN;
        } else {
            instanceConfig.watchedEpisodesCssClass = watchedStateClasses.FAINT;
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
            let linkText;
            let hideWatched;
            if (instanceConfig.watchedEpisodesCssClass === watchedStateClasses.FAINT) {
                instanceConfig.watchedEpisodesCssClass = watchedStateClasses.HIDDEN;
                linkText = "SHOW WATCHED";
                hideWatched = true;
            } else {
                instanceConfig.watchedEpisodesCssClass = watchedStateClasses.FAINT;
                linkText = "HIDE WATCHED";
                hideWatched = false;
            }
            let config = JSON.parse(
                localStorage.getItem(localStorageKeys.CONFIG) || "{}"
            );
            config.hideWatched = hideWatched;
            localStorage.setItem(localStorageKeys.CONFIG, JSON.stringify(config))

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
        });
    };

    $(document).ready(function() {
        loadConfig();
        $('#show-filter-select').select2({
          placeholder: 'Select shows to exclude...',
          allowClear: true,
          closeOnSelect: false,
          width: "100%",
        });
        registerListeners();
        setWatchedDisplayState(instanceConfig.watchedEpisodesCssClass);

        var colourSetting = Cookies.get('colour');
        if (colourSetting === undefined) {
            Cookies.set('colour', '1');
        } else if (colourSetting === '0') {
            disableColours();
        }
    });
})(jQuery);
