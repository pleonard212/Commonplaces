(function() {
    'use strict';

	angular
		.module('DiggingApp')
		.controller('MainController', MainController);

		function MainController($scope, $routeParams, $route, $log, $http, $location, webConfig) {
            var vm = this;

            vm.hideLanding = true;
            vm.hideSearchForm = false;
            vm.landingPage = true;
            vm.webConfig = webConfig;
            vm.description = false;
            vm.formData = {};
            vm.formData.duplicates = "ignore";
            vm.formData.sorting = 1;
            vm.dbActive = vm.webConfig.databases[0].dbname;

            vm.openLandingPage = function() {
                if (vm.hideSearchForm) {
                    vm.hideSearchForm = false;
                }
            }
		}
})();
