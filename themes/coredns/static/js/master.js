$(document).ready(function() {
    $(':header[id]').each(function() {
        var anchor = document.createElement('a')
        anchor.href = '#' + this.id
        $(this).wrapInner(anchor)
    });
});
