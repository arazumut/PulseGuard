$(function () {
    $('.knob').knob({
        draw: function() {
            this.i.css({
                'margin-right': '-' + ((this.w * 3 / 4 + 2) >> 0) + 'px',
                'margin-left': 'auto'
            });
        }
    });
});