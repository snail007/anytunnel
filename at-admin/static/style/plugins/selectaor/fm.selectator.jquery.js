/*
 Selectator jQuery Plugin
 A plugin for select elements
 version 2.2, Dac 11th, 2015
 by Ingi P. Jacobsen

 The MIT License (MIT)

 Copyright (c) 2013 Faroe Media

 Permission is hereby granted, free of charge, to any person obtaining a copy of
 this software and associated documentation files (the "Software"), to deal in
 the Software without restriction, including without limitation the rights to
 use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 the Software, and to permit persons to whom the Software is furnished to do so,
 subject to the following conditions:

 The above copyright notice and this permission notice shall be included in all
 copies or substantial portions of the Software.

 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

(function($) {
	$.selectator = function (element, options) {
		var defaults = {
			prefix: 'selectator_',
			height: 'auto',
			useDimmer: false,
			useSearch: true,
			showAllOptionsOnFocus: false,
			selectFirstOptionOnSearch: true,
			keepOpen: false,
			searchCallback: function(){},
			labels: {
				search: 'Search...'
			}
		};

		var plugin = this;
		plugin.settings = {};
		var $source_element = $(element);
		var $container_element = null;
		var $chosenitems_element = null;
		var $input_element = null;
		var $textlength_element = null;
		var $options_element = null;
		var is_single = $source_element.attr('multiple') === undefined;
		var is_multiple = !is_single;
		var has_visible_options = true;
		var key = {
			backspace: 8,
			tab:       9,
			enter:    13,
			shift:    16,
			ctrl:     17,
			alt:      18,
			capslock: 20,
			escape:   27,
			pageup:   33,
			pagedown: 34,
			end:      35,
			home:     36,
			left:     37,
			up:       38,
			right:    39,
			down:     40
		};

		
		
		// INITIALIZE PLUGIN
		plugin.init = function () {
			plugin.settings = $.extend({}, defaults, options);

			//// ================== CREATE ELEMENTS ================== ////
			// dimmer
			if (plugin.settings.useDimmer) {
				if ($('#' + plugin.settings.prefix + 'dimmer').length === 0) {
					var $dimmer_element = $(document.createElement('div'));
					$dimmer_element.attr('id', plugin.settings.prefix + 'dimmer');
					$dimmer_element.hide();
					$(document.body).prepend($dimmer_element);
				}
			}
			// source element
			$source_element.addClass('selectator');
			// container element
			$container_element = $(document.createElement('div'));
			if ($source_element.attr('id') !== undefined) {
				$container_element.attr('id', plugin.settings.prefix + $source_element.attr('id'));
			}
			$container_element.addClass(plugin.settings.prefix + 'element ' + (is_multiple ? 'multiple ' : 'single ') + 'options-hidden');
			if (!plugin.settings.useSearch) {
				$container_element.addClass('disable_search');
			}
			$container_element.css({
				width: $source_element.css('width'),
				minHeight: $source_element.css('height'),
				padding: $source_element.css('padding'),
				'flex-grow': $source_element.css('flex-grow'),
				position: 'relative'
			});
			if (plugin.settings.height === 'element') {
				$container_element.css({
					height: $source_element.outerHeight() + 'px'
				});
			}
			// textlength element
			$textlength_element = $(document.createElement('span'));
			$textlength_element.addClass(plugin.settings.prefix + 'textlength');
			$textlength_element.css({
				position: 'absolute',
				visibility: 'hidden'
			});
			$container_element.append($textlength_element);
			// chosen items element
			$chosenitems_element = $(document.createElement('div'));
			$chosenitems_element.addClass(plugin.settings.prefix + 'chosen_items');
			$container_element.append($chosenitems_element);
			// input element
			$input_element = $(document.createElement('input'));
			$input_element.addClass(plugin.settings.prefix + 'input');
			$input_element.attr('tabindex', $source_element.attr('tabindex'));
			if (!plugin.settings.useSearch) {
				$input_element.attr('readonly', true);
				$input_element.css({
					'width': '0px',
					'height': '0px',
					'overflow': 'hidden',
					'border': 0,
					'padding': 0,
					'position': 'absolute'
				});
			} else {
				if (is_single) {
					$input_element.attr('placeholder', plugin.settings.labels.search);
				} else {
					$input_element.width(20);
				}
			}
			$input_element.attr('autocomplete', 'false');
			$container_element.append($input_element);
			// options element
			$options_element = $(document.createElement('ul'));
			$options_element.addClass(plugin.settings.prefix + 'options');

			$container_element.append($options_element);
			$source_element.after($container_element);
			$source_element.hide();

			// Add scrollator if found
			if (typeof Scrollator !== 'undefined') {
				$options_element.scrollator({
					zIndex: 1001,
					customClass : 'ease_preventOverlay'
				});
			}


			//// ================== BIND ELEMENTS EVENTS ================== ////
			// source element
			$source_element.change(function () {
				regenerateChosenItems();
			});
			// container element
			$container_element.on('focus', function (e) {
				$input_element.focus();
				$input_element.trigger('focus');
			});
			$container_element.on('mousedown', function (e) {
				e.preventDefault();
				$input_element.focus();
				$input_element.trigger('focus');
				// put text caret to end of search field
				if ($input_element[0].setSelectionRange) {
					$input_element[0].setSelectionRange($input_element.val().length, $input_element.val().length);
				} else if ($input_element[0].createTextRange) {
					var range = $input_element[0].createTextRange();
					range.collapse(true);
					range.moveEnd('character', $input_element.val().length);
					range.moveStart('character', $input_element.val().length);
					range.select();
				}
			});
			$container_element.on('mouseup', function (e) {
			});
			$container_element.on('click', function (e) {
				$input_element.focus();
				$input_element.trigger('focus');
			});
			$container_element.on('dblclick', function (e) {
				$input_element.select();
				$input_element.trigger('focus');
			});
			// input element
			$input_element.on('click', function (e) {
			});
			$input_element.on('dblclick', function (e) {
			});
			$input_element.on('keydown', function (e) {
				var keyCode = e.keyCode || e.which;
				var $active = null;
				var $newActive = null;
				switch (keyCode) {
					case key.up:
						e.preventDefault();
						showDropdown();
						$active = $options_element.find('.active');
						if ($active.length !== 0) {
							$newActive = $active.prevUntil('.' + plugin.settings.prefix + 'option:visible').add($active).first().prev('.' + plugin.settings.prefix + 'option:visible');
							$active.removeClass('active');
							$newActive.addClass('active');
						} else {
							$options_element.find('.' + plugin.settings.prefix + 'option').filter(':visible').last().addClass('active');
						}
						scrollToActiveOption();
						break;
					case key.down:
						e.preventDefault();
						showDropdown();
						$active = $options_element.find('.active');
						if ($active.length !== 0) {
							$newActive = $active.nextUntil('.' + plugin.settings.prefix + 'option:visible').add($active).last().next('.' + plugin.settings.prefix + 'option:visible');
							$active.removeClass('active');
							$newActive.addClass('active');
						} else {
							$options_element.find('.' + plugin.settings.prefix + 'option').filter(':visible').first().addClass('active');
						}
						scrollToActiveOption();
						break;
					case key.escape:
						e.preventDefault();
						break;
					case key.enter:
						e.preventDefault();
						$active = $options_element.find('.active');
						if ($active.length !== 0) {
							selectOption();
						} else {
							if ($input_element.val() !== '') {
								plugin.settings.searchCallback($input_element.val());
							}
						}
						resizeSearchInput();
						break;
					case key.tab:
						e.preventDefault();
						$active = $options_element.find('.active');
						if ($active.length !== 0) {
							selectOption();
						} else {
							if ($input_element.val() !== '') {
								plugin.settings.searchCallback($input_element.val());
							}
						}
						resizeSearchInput();
						break;
					case key.backspace:
						if (plugin.settings.useSearch) {
							if ($input_element.val() === '' && is_multiple) {
								$source_element.find('option:selected').last()[0].selected = false;
								$source_element.trigger('change');
								regenerateChosenItems();
							}
							resizeSearchInput();
						} else {
							e.preventDefault();
						}
						break;
					default:
						resizeSearchInput();
						break;
				}
			});
			$input_element.on('keyup', function (e) {
				e.preventDefault();
				e.stopPropagation();
				var keyCode = e.which;
				switch (keyCode) {
					case key.escape:
						hideDropdown();
						break;
					case key.enter:
						if (!plugin.settings.keepOpen) {
							hideDropdown();
						}
						break;
					case key.left:
					case key.right:
					case key.up:
					case key.down:
					case key.tab:
					case key.shift:
						// Prevent any action
						break;
					default:
						search();
						break;
				}
				if ($container_element.hasClass('options-hidden') && (keyCode === key.left || keyCode === key.right || keyCode === key.up || keyCode === key.down)) {
					showDropdown();
				}
				resizeSearchInput();
			});
			$input_element.on('focus', function (e) {
				$container_element.addClass('focused');
				if (is_single || plugin.settings.showAllOptionsOnFocus || !plugin.settings.useSearch) {
					showDropdown();
				}
			});
			$input_element.on('blur', function (e) {
				$container_element.removeClass('focused');
				hideDropdown();
			});
			

			// bind option events
			$container_element.delegate('.' + plugin.settings.prefix + 'option', 'mouseover', function (e) {
				var $active = $options_element.find('.active');
				$active.removeClass('active');
				var $this = $(this);
				$this.addClass('active');
			});
			$container_element.delegate('.' + plugin.settings.prefix + 'option', 'mousedown', function (e) {
				e.preventDefault();
				e.stopPropagation();
			});
			$container_element.delegate('.' + plugin.settings.prefix + 'option', 'mouseup', function (e) {
				selectOption();
			});
			$container_element.delegate('.' + plugin.settings.prefix + 'option', 'click', function (e) {
				e.stopPropagation();
			});

			regenerateOptions();
			regenerateChosenItems();
		};


		// RESIZE INPUT
		var resizeSearchInput = function () {
			$textlength_element.text($input_element.val());
			if (is_multiple) {
				var width = $textlength_element.width() > ($container_element.width() - 30) ? ($container_element.width() - 30) : ($textlength_element.width() + 30);
				$input_element.css({ width: width + 'px' });
			}
		};


		// REGENERATE CHOSEN ITEMS
		var regenerateChosenItems = function () {
			$chosenitems_element.empty();
			$source_element.find('option').each(function () {
				var $option = $(this);
				if (this.selected) {
					var $item_element = $(document.createElement('div'));
					$item_element.addClass(plugin.settings.prefix + 'chosen_item');
					$item_element.addClass(plugin.settings.prefix + 'value_' + $option.val().replace(/\W/g, ''));

					// class
					if ($option.attr('class') !== undefined) {
						$item_element.addClass($option.attr('class'));
					}
					// left
					if ($option.data('left') !== undefined) {
						var $left_element = $(document.createElement('div'));
						$left_element.addClass(plugin.settings.prefix + 'chosen_item_left').html($option.attr('data-left'));
						$item_element.append($left_element);
					}
					// right
					if ($option.data('right') !== undefined) {
						var $right_element = $(document.createElement('div'));
						$right_element.addClass(plugin.settings.prefix + 'chosen_item_right').html($option.attr('data-right'));
						$item_element.append($right_element);
					}
					// title
					var $title_element = $(document.createElement('div'));
					$title_element.addClass(plugin.settings.prefix + 'chosen_item_title').html($option.html());
					$item_element.append($title_element);
					// subtitle
					if ($option.data('subtitle') !== undefined) {
						var $subtitle_element = $(document.createElement('div'));
						$subtitle_element.addClass(plugin.settings.prefix + 'chosen_item_subtitle').html($option.attr('data-subtitle'));
						$item_element.append($subtitle_element);
					}
					// remove button
					var $button_remove_element = $(document.createElement('div'));
					$button_remove_element.data('source_option_element', this);
					$button_remove_element.addClass(plugin.settings.prefix + 'chosen_item_remove');
					$button_remove_element.on('mousedown', function (e) {
					});
					$button_remove_element.on('mouseup', function (e) {
						$(this).data('source_option_element').selected = false;
						$source_element.trigger('change');
						search();
						regenerateChosenItems();
					});
					$button_remove_element.html('X');
					$item_element.append($button_remove_element);
					$chosenitems_element.append($item_element);
				}
			});
		};


		// REGENERATE OPTIONS
		var regenerateOptions = function () {
			$options_element.empty();
			var optionsArray = [];
			$source_element.children().each(function () {
				if ($(this).prop('tagName').toLowerCase() === 'optgroup') {
					var $group = $(this);
					if ($group.children('option').length !== 0) {
						var groupOptionsArray = [];
						$group.children('option').each(function () {
							groupOptionsArray.push({
								type: 'option',
								text: $(this).html(),
								element: this
							});
						});
						optionsArray.push({
							type: 'group',
							text: $group.attr('label'),
							options: groupOptionsArray,
							element: $group
						});
					}
				} else {
					optionsArray.push({
						type: 'option',
						text: $(this).html(),
						element: this
					});
				}
			});

			$(optionsArray).each(function () {
				if (this.type === 'group') {
					var $group_element = $(document.createElement('li'));
					$group_element.addClass(plugin.settings.prefix + 'group');
					if ($(this.element).attr('class') !== undefined) {
						$group_element.addClass($(this.element).attr('class'));
					}
					$group_element.html($(this.element).attr('label'));
					$options_element.append($group_element);

					$(this.options).each(function () {
						var option = createOption.call(this.element, true);
						$options_element.append(option);
					});

				} else {
					var option = createOption.call(this.element, false);
					$options_element.append(option);
				}
			});
			search();
		};
		

		// CREATE RESULT OPTION
		var createOption = function (isGroupOption) {
			// holder li
			var $option = $(document.createElement('li'));
			$option.data('source_option_element', this);
			$option.addClass(plugin.settings.prefix + 'option');
			$option.addClass(plugin.settings.prefix + 'value_' + $(this).val().replace(/\W/g, ''));
			if (isGroupOption) {
				$option.addClass(plugin.settings.prefix + 'group_option');
			}
			if (this.selected) {
				$option.addClass('active');
			}
			// class
			if ($(this).attr('class') !== undefined) {
				$option.addClass($(this).attr('class'));
			}
			// left
			if ($(this).data('left') !== undefined) {
				var $left_element = $(document.createElement('div'));
				$left_element.addClass(plugin.settings.prefix + 'option_left').html($(this).attr('data-left'));
				$option.append($left_element);
			}
			// right
			if ($(this).data('right') !== undefined) {
				var $right_element = $(document.createElement('div'));
				$right_element.addClass(plugin.settings.prefix + 'option_right').html($(this).attr('data-right'));
				$option.append($right_element);
			}
			// title
			var $title_element = $(document.createElement('div'));
			$title_element.addClass(plugin.settings.prefix + 'option_title').html($(this).html());
			$option.append($title_element);
			// subtitle
			if ($(this).data('subtitle') !== undefined) {
				var $subtitle_element = $(document.createElement('div'));
				$subtitle_element.addClass(plugin.settings.prefix + 'option_subtitle').html($(this).attr('data-subtitle'));
				$option.append($subtitle_element);
			}

			return $option;
		};
		
		
		// SEARCH
		var search = function () {
			// bool true if search field is considered empty
			var searchIsEmpty = $input_element.val().replace(/\s/g, '') === '';
			// bool true if any options are visible
			has_visible_options = false;
			// get sanitized search text
			var searchFor = $input_element.val().toLowerCase();
			// iterate through the options
			$options_element.find('.' + plugin.settings.prefix + 'option').each(function () {
				var $this = $(this);
				var source_option_element = $this.data('source_option_element');
				// show if:
				// (item is not selected  or  if single select)
				// and 
				//     use search
				//         and search is empty  or  text matches the input box
				//     or not using search
				if (
					(!source_option_element.selected || is_single) 
					&& (
						plugin.settings.useSearch
						&& (
							searchIsEmpty 
							|| $(source_option_element).html().toLowerCase().indexOf(searchFor) !== -1
						)
						|| !plugin.settings.useSearch
					)
				) {
					$this.show();
					has_visible_options = !has_visible_options ? true : has_visible_options;
				} else {
					$this.hide();
				}
			});
			// iterate through the groups
			$options_element.find('.' + plugin.settings.prefix + 'group').each(function () {
				var $this = $(this);
				var has_visible_options = false;
				$this.nextUntil('.' + plugin.settings.prefix + 'group').each(function () {
					var $option = $(this);
					if ($option.css('display') != 'none') {
						has_visible_options = true;
						return false;
					}
				});
				// show if the group has any visible children
				if (has_visible_options) {
					$this.show();
				} else {
					$this.hide();
				}
			});
			showDropdown();
			if (is_multiple) {
				$options_element.find('.active').removeClass('active');
				if (!searchIsEmpty) {
					$options_element.find('.' + plugin.settings.prefix + 'option').filter(':visible').first().addClass('active');
				}
			}
		};


		// SHOW OPTIONS AND DIMMER
		var showDropdown = function () {
			if ($input_element.is(':focus') && (has_visible_options || is_single )) {
				$container_element.removeClass('options-hidden').addClass('options-visible');
				if (plugin.settings.useDimmer) {
					$('#' + plugin.settings.prefix + 'dimmer').show();
				}
				setTimeout(function () {
					$options_element.css('top', ($container_element.outerHeight() + (is_multiple ? 0 : $input_element.outerHeight()) - 1) + 'px');
					if (typeof Scrollator !== 'undefined') {
						$options_element.data('scrollator').show();
					}
				}, 1);
				scrollToActiveOption();
			} else {
				hideDropdown();
			}
		};


		// HIDE OPTIONS AND DIMMER
		var hideDropdown = function () {
			$container_element.removeClass('options-visible').addClass('options-hidden');
			if (typeof Scrollator !== 'undefined') {
				$options_element.data('scrollator').hide();
			}
			if (plugin.settings.useDimmer) {
				$('#' + plugin.settings.prefix + 'dimmer').hide();
			}
		};


		// SCROLL TO ACTIVE OPTION IN OPTIONS LIST
		var scrollToActiveOption = function () {
			var $active_element = $options_element.find('.' + plugin.settings.prefix + 'option.active');
			if ($active_element.length > 0) {
				$options_element.scrollTop($options_element.scrollTop() + $active_element.position().top - $options_element.height()/2 + $active_element.height()/2);
			}
		};


		// SELECT ACTIVE OPTION
		var selectOption = function () {
			// select option
			var $active = $options_element.find('.active');
			$active.data('source_option_element').selected = true;
			$source_element.trigger('change');
			$input_element.val('');
			search();
			regenerateChosenItems();
			if (!plugin.settings.keepOpen) {
				hideDropdown();
			}
		};


		// REFRESH PLUGIN
		plugin.refresh = function () {
			regenerateChosenItems();
		};


		// REMOVE PLUGIN AND REVERT SELECT ELEMENT TO ORIGINAL STATE
		plugin.destroy = function () {
			$container_element.remove();
			$.removeData(element, 'selectator');
			$source_element.show();
			if ($('.' + plugin.settings.prefix + 'element').length === 0) {
				$('#' + plugin.settings.prefix + 'dimmer').remove();
			}
		};

		
		// Initialize plugin
		plugin.init();
	};
	
	$.fn.selectator = function(options) {
		options = options !== undefined ? options : {};
		return this.each(function () {
			if (typeof(options) === 'object') {
				if (undefined === $(this).data('selectator')) {
					var plugin = new $.selectator(this, options);
					$(this).data('selectator', plugin);
				}
			} else if ($(this).data('selectator')[options]) {
				$(this).data('selectator')[options].apply(this, Array.prototype.slice.call(arguments, 1));
			} else {
				$.error('Method ' + options + ' does not exist in $.selectator');
			}
		});
	};

}(jQuery));

$(function () {
	$('.selectator').each(function () {
		var $this = $(this);
		var options = {};
		$.each($this.data(), function (key, value) {
			if (key.substring(0, 10) == 'selectator') {
				options[key.substring(10, 11).toLowerCase() + key.substring(11)] = value;
			}
		});
		$this.selectator(options);
	});
});