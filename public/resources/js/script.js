var autocpl;

$(function()
{
	var KEY_DOWN = 40,
		KEY_UP = 38,
		KEY_TAB = 9,
		KEY_ENTER = 13;
	autocpl = $('.autocomplete');

	//$('li').tipsy({live: true});	// TODO.

	$('#search').focus();

	$('form.host-search').submit(function(event)
	{
		var host = $.trim($('#search').val());
		
		if (!host)
			return suppress(event);

		window.location = '/report/' + encodeURIComponent(host);
		return suppress(event);
	});

	autocpl.on('click', '.autocpl-suggestion', function()
	{
		useAutocplSuggestion($(this));
		$('form.host-search').submit();
	});

	$('#search').keyup(function(event)
	{
		if (event.keyCode == KEY_UP || event.keyCode == KEY_DOWN
			|| event.keyCode == KEY_TAB || event.keyCode == KEY_ENTER)
			return suppress(event);

		var currentChoice = $('.autocpl-suggestion.sel').first();

		var input = $(this).val().trim();
		if (!input)
		{
			autocpl.hide();
			return;
		}

		$.getJSON('/app/suggest/'+encodeURIComponent(input), function(data)
		{
			autocpl.empty().show();
			moveAutocomplete();

			if (data.length > 0)
			{
				for (var i in data)
				{
					var inputStart = data[i].indexOf(input);
					var before = data[i].slice(0, inputStart);
					var matched = data[i].slice(inputStart, inputStart + input.length);
					var after = data[i].slice(inputStart + input.length);
					sugg = before + '<span class="matched">' + matched + '</span>' + after;
					autocpl.append('<div class="autocpl-suggestion">'+sugg+'</div>');
				}
			}
			else
			{
				autocpl.append('<div class="autocpl-nosugg">Not found in our database</div>');
			}
		});
	}).keydown(function(event)
	{
		var currentChoice = $('.autocpl-suggestion.sel').first();
		var choiceSelectionIsNew = false;

		if (event.keyCode == KEY_ENTER)
		{
			if ($('.autocpl-suggestion.sel:visible').length > 0)
			{
				useAutocplSuggestion(currentChoice);
				return;
			}
		}
		else if (event.keyCode == KEY_TAB)
		{
			if (currentChoice.length > 0)
			{
				useAutocplSuggestion(currentChoice);
				return suppress(event);
			}
			else
			{
				autocpl.hide();
				return;
			}
		}
		else if (event.keyCode == KEY_DOWN)
		{
			if (!currentChoice.hasClass('autocpl-suggestion'))
			{
				currentChoice = $('.autocpl-suggestion').first().mouseover();
				choiceSelectionIsNew = true;
			}

			if (!choiceSelectionIsNew)
			{
				currentChoice.removeClass('sel');
				currentChoice.next('.autocpl-suggestion').mouseover();
			}

			moveCursorToEnd(this);
			return;
		}
		else if (event.keyCode == KEY_UP)
		{
			if (!currentChoice.hasClass('autocpl-suggestion'))
			{
				currentChoice = $('.autocpl-suggestion').last().mouseover();
				choiceSelectionIsNew = true;
			}

			if (!choiceSelectionIsNew)
			{
				currentChoice.removeClass('sel');
				currentChoice.prev('.autocpl-suggestion').mouseover();
			}

			moveCursorToEnd(this);
			return;
		}
	});

	autocpl.on('mouseover', '.autocpl-suggestion', function() {
		$(this).addClass('sel');
	}).on('mouseleave', '.autocpl-suggestion', function() {
		$(this).removeClass('sel');
	});

	$(window).resize(moveAutocomplete);

});

function useAutocplSuggestion(choice)
{
	$('#search').val(choice.text());
	autocpl.hide();
}

function moveAutocomplete()
{
	var search = $('#search');
	var pos = search.position();
	var bottom = pos.top + search.outerHeight() + 2;
	var left = pos.left;
	var width = search.width() - 1;

	$('.autocomplete').css({
		'top': bottom + 'px',
		'left': left + 'px',
		'width': width + 'px'
	});
}


function moveCursorToEnd(el)	// Courtesy of http://css-tricks.com/snippets/javascript/move-cursor-to-end-of-input/
{
	if (typeof el.selectionStart == "number")
		el.selectionStart = el.selectionEnd = el.value.length;
	else if (typeof el.createTextRange != "undefined")
	{
		el.focus();
		var range = el.createTextRange();
		range.collapse(false);
		range.select();
	}
}



function suppress(event)
{
	if (!event)
		return false;
	if (event.preventDefault)
		event.preventDefault();
	if (event.stopPropagation)
		event.stopPropagation();
	if (event.cancelBubble)
		event.cancelBubble = true;
	return false;
}