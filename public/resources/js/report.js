var parts = window.location.pathname.split('/');
var host = parts[parts.length - 1];
$('title').text("Report of errors and typos on "+host);

var RPT_NEW = 0,
	RPT_ACK = 1,
	RPT_FIX = 2,
	RPT_REJ = 4,
	RPT_DEL = 8;

var data = [];
var reportGroups = {};	// suggestions grouped together, keyed by an arbitrary report ID for easy retrieval


$(function()
{
	$('.host').text(host);

	$.getJSON('/app/report/'+host).done(function(pages, status, jqxhr)
	{
		$('.host-found').show();
		$('.host-not-found, .loading').hide();

		// Turn the objects into arrays so the templater can render it
		var pagesArray = [];
		var suggCount = 0;	// number of suggestions on open issues
		var id = 0;		// a rendering ID used for client-side convenience
		for (var pageUri in pages)
		{
			pagesArray.push({
				link: "http://" + host + pageUri,
				uri: pageUri,
				reports: []
			});

			var reports = pages[pageUri];

			var reportCount = 0;
			for (var report in reports)
			{
				if (!shouldShowReport(reports[report].Status))
					continue;

				pagesArray[pagesArray.length - 1].reports.push({
					originalText: report,
					status: reports[report].Status,
					suggestions: reports[report].Suggestions,
					_id: id
				});
				reportCount++;

				suggCount += reports[report].Suggestions.length;
				
				// for convenience (and by necessity) for when user acts on a report
				reports[report].host = host;
				reports[report].uri = pageUri;
				reportGroups[id] = reports[report];

				id++;	// increment the rendering ID used just client-side
			}

			if (reportCount == 0)
				pagesArray.pop();	// no need showing a page with no open issues
		}

		// Show suggestion count
		if (suggCount == 0)
			$('.report-count').html(suggCount + " outstanding suggestions");
		else
			$('.report-count').html(suggCount+' suggestion'+(suggCount == 1 ? "" : "s"));

		// Save and render the data
		data = pagesArray;
		$('#reports').html(render('tpl-reports', pagesArray));

	}).fail(function(jqxhr, err, msg)
	{
		$('.host-found, .loading').hide();

		if (jqxhr.status == 404) // hostname not found
			$('.host-not-found').show();
	});


	////////////// Report actions ///////////////
	
	$('#reports').on('click', '.controls li.ack', function() {
		var renderingID = $(this).closest('tr.next-sugg').data('id');
		$(this).add('tr.report-'+renderingID)
			.removeClass('ack fixed reject delete').addClass('ack');
		$.post('/app/update/' + reportGroups[renderingID].ID + '/ack', {
			host: reportGroups[renderingID].host,
			uri: reportGroups[renderingID].uri
		});
	}).on('click', '.controls li.fixed', function() {
		if (!confirm("Are you sure this is fixed? This report will disappear."))
			return;
		var report = $(this).closest('tr.next-sugg');
		var renderingID = report.data('id');
		$(this).add('tr.report-'+renderingID)
			.removeClass('ack fixed reject delete').addClass('fixed');
		$.post('/app/update/' + reportGroups[renderingID].ID + '/fixed', {
			host: reportGroups[renderingID].host,
			uri: reportGroups[renderingID].uri
		}, function() {
			$('.report-'+renderingID).fadeOut();
		});
	})/*
	.on('click', '.controls li.reject', function()
	{
		// Currently not enabled because, right now, reject
		// seems redundant to delete, but in the future they
		// may mean different things.

		if (!confirm("Really reject this report? It will disappear."))
			return;
		var renderingID = $(this).closest('tr.next-sugg').data('id');
		$(this).add('tr.report-'+renderingID)
			.removeClass('ack fixed reject delete').addClass('reject');
		$.post('/app/update/' + reportGroups[renderingID].ID + '/reject', {
			host: reportGroups[renderingID].host,
			uri: reportGroups[renderingID].uri
		});
		report.fadeOut();
	})*/.on('click', '.controls li.delete', function() {
		if (!confirm("Are you sure you want to delete this?"))
			return;
		var renderingID = $(this).closest('tr.next-sugg').data('id');
		$(this).add('tr.report-'+renderingID)
			.removeClass('ack fixed reject delete').addClass('delete');
		$.post('/app/update/' + reportGroups[renderingID].ID + '/delete', {
			host: reportGroups[renderingID].host,
			uri: reportGroups[renderingID].uri
		});
		$('.report-'+renderingID).fadeOut();
	});
});

function render(templateID, context)
{
	var tpl = $('#' + templateID).text();
	return Mark.up(tpl, context);
}

function shouldShowReport(status)
{
	return (status & RPT_DEL) == 0
		&& (status & RPT_REJ) == 0
		&& (status & RPT_FIX) == 0;
}

Mark.pipes.flagsToClass = function(flags)
{
	flags = parseInt(flags);
	if (flags & RPT_DEL)
		return "delete";
	if (flags & RPT_REJ)
		return "reject";
	if (flags & RPT_FIX)
		return "fixed";
	if (flags & RPT_ACK)
		return "ack";
	return "";
};

Mark.pipes.hasnewline = function(str)
{
	return str.indexOf('\n') > -1;
};