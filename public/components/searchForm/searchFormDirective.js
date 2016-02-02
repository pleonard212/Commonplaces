(function() {
    'use strict';
    angular
        .module('DiggingApp')
        .directive('searchForm', searchForm);

    function searchForm($location, $routeParams, $log, $http, URL, sortKeys) {
        var hideLandingPage = function() {
            angular.element('.hiding-element').each(function() {
                hideElement(angular.element(this));
            })
        }
        var hideElement = function(element) {
            if (element instanceof jQuery === false) {
                element = angular.element(element.currentTarget).siblings(".hiding-element");
            }
            element.velocity('slideUp', {
                duration: 250,
                easing: "easeOut"
            });
            element.siblings(".close").addClass("closed");
            element.parent().velocity({
                "padding-top": 0,
                "padding-bottom": 0,
                "margin-bottom": "5px",
                "margin-top": 0
            }, {
                duration: 250,
                easing: "easeOut",
                queue: false
            });
            var titleBar = element.parent().find('h4');
            titleBar.velocity({
                'font-size': "100%",
                'color': '#155F83',
                "margin-bottom": '0px',
                "padding": "5px"
            }, {
                duration: 250,
                easing: "easeOut",
                queue: false
            });
            titleBar.css('cursor', 'pointer');
            titleBar.off().on('click touchstart', function() {
                showElement(angular.element(this).parent().find('.hiding-element'));
            });
            element.parent().find('.glyphicon-chevron-right').show();
            element.parent().find('.glyphicon-chevron-down').hide();
        }
        var showElement = function(element) {
            element
                .velocity('slideDown', {
                    duration: 250,
                    easing: "easeOut"
                });
            element.siblings(".close").removeClass("closed");
            element.parent().velocity({
                "padding": "15px",
                "margin-bottom": "15px",
                "margin-top": "15px"
            });
            var titleBar = element.parent().find('h4');
            titleBar.off();
            titleBar.velocity({
                'font-size': "120%",
                'color': 'black',
                "margin-bottom": '15px',
                "padding": 0
            });
            titleBar.on("click touchstart", function() {
                hideElement(element);
            });
            element.parent().find('.glyphicon-chevron-right').hide();
            element.parent().find('.glyphicon-chevron-down').show();
        }
        var showLandingPage = function() {
            angular.element("#landing-page-container").find('.hiding-element').each(function() {
                showElement(angular.element(this));
            })
        }
        var sortObject = function(obj) {
            var arr = [];
            var prop;
            for (prop in obj) {
                if (obj.hasOwnProperty(prop)) {
                    arr.push({
                        'key': prop,
                        'value': obj[prop]
                    });
                }
            }
            arr.sort(function(a, b) {
                return b.value - a.value;
            });
            return arr; // returns array
        }
        return {
            templateUrl: 'components/searchForm/searchForm.html',
            link: function(scope) {
                angular.element('[data-toggle="popover"]').popover();
                scope.sorting = "No Sorting";
                scope.submit = function() {
                    scope.results = [];
                    if (scope.main.formData.duplicates != "ignore") {
                        delete scope.main.formData.duplicates;
                    }
                    var urlString = URL.objectToString(scope.main.formData);
                    $log.debug(urlString)
                    if (urlString.length === 0) {
                        alert("You haven't searched for anything, please fill in one of the search boxes");
                    } else {
                        hideLandingPage();
                        $location.url('/query/' + scope.main.dbActive + '/search?' + urlString)
                    }
                };
                scope.commonplaceTerms = {};
                scope.commonplaceSubmit = function() {
                    hideLandingPage();
                    var queryString = URL.objectToString(scope.main.commonplace)
                    var urlString = "/commonplace/" + scope.main.dbActive + "/search?" + queryString;
                    $log.debug(urlString);
                    $location.url(urlString)
                }
                scope.toggleForm = function() {
                    if (scope.hideForm) {
                        angular.element('.hiding-element').velocity('slideDown');
                    } else {
                        hideLandingPage();
                    }
                }
                scope.hideElement = hideElement;
                scope.selectSorting = function(sortId) {
                    scope.main.formData.sorting = sortId;
                    scope.sorting = sortKeys.keys[sortId].label;
                }
                scope.bibleFilter = "No filter";
                scope.bibleFiltering = function(filtering) {
                    if (filtering === 0) {
                        scope.main.formData.bible = "ignore";
                        scope.bibleFilter = "Filter out Bible sources";
                    } else if (filtering == 1) {
                        scope.main.formData.bible = "only";
                        scope.bibleFilter = "Filter out non-Bible sources";
                    } else if (filtering === 2) {
                        scope.main.formData.bible = "all";
                        scope.bibleFilter = "off";
                    }
                    $log.debug(scope.main.formData.bible)
                }
                scope.showLatinAuthorList = false;
                scope.listAuthors = function() {
                    if (!scope.showLatinAuthorList) {
                        var urlString = "/api/getLatinAuthors";
                        $http.get(urlString).then(function(response) {
                            scope.latinAuthors = sortObject(response.data);
                            $log.debug(scope.latinAuthors)
                            scope.showLatinAuthorList = true;
                        });
                    } else {
                        scope.showLatinAuthorList = false;
                    }
                }
                scope.$watch("main.hideSearchForm", function(currentValue) {
                    if (currentValue) {
                        hideLandingPage();
                    } else {
                        showLandingPage();
                    }
                });
            }
        }
    }
})();
