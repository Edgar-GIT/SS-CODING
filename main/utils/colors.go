package utils

type Color string

const (
    // Reset
    Reset Color = "\033[0m"

    // Regular Colors
    Black   Color = "\033[0;30m"
    Red     Color = "\033[0;31m"
    Green   Color = "\033[0;32m"
    Yellow  Color = "\033[0;33m"
    Blue    Color = "\033[0;34m"
    Purple  Color = "\033[0;35m"
    Cyan    Color = "\033[0;36m"
    White   Color = "\033[0;37m"

    // Bold
    BoldBlack   Color = "\033[1;30m"
    BoldRed     Color = "\033[1;31m"
    BoldGreen   Color = "\033[1;32m"
    BoldYellow  Color = "\033[1;33m"
    BoldBlue    Color = "\033[1;34m"
    BoldPurple  Color = "\033[1;35m"
    BoldCyan    Color = "\033[1;36m"
    BoldWhite   Color = "\033[1;37m"

    // Underline
    UnderlineBlack   Color = "\033[4;30m"
    UnderlineRed     Color = "\033[4;31m"
    UnderlineGreen   Color = "\033[4;32m"
    UnderlineYellow  Color = "\033[4;33m"
    UnderlineBlue    Color = "\033[4;34m"
    UnderlinePurple  Color = "\033[4;35m"
    UnderlineCyan    Color = "\033[4;36m"
    UnderlineWhite   Color = "\033[4;37m"

    // Background
    BgBlack   Color = "\033[40m"
    BgRed     Color = "\033[41m"
    BgGreen   Color = "\033[42m"
    BgYellow  Color = "\033[43m"
    BgBlue    Color = "\033[44m"
    BgPurple  Color = "\033[45m"
    BgCyan    Color = "\033[46m"
    BgWhite   Color = "\033[47m"

    // High Intensity
    HiBlack   Color = "\033[0;90m"
    HiRed     Color = "\033[0;91m"
    HiGreen   Color = "\033[0;92m"
    HiYellow  Color = "\033[0;93m"
    HiBlue    Color = "\033[0;94m"
    HiPurple  Color = "\033[0;95m"
    HiCyan    Color = "\033[0;96m"
    HiWhite   Color = "\033[0;97m"

    // Bold High Intensity
    BoldHiBlack   Color = "\033[1;90m"
    BoldHiRed     Color = "\033[1;91m"
    BoldHiGreen   Color = "\033[1;92m"
    BoldHiYellow  Color = "\033[1;93m"
    BoldHiBlue    Color = "\033[1;94m"
    BoldHiPurple  Color = "\033[1;95m"
    BoldHiCyan    Color = "\033[1;96m"
    BoldHiWhite   Color = "\033[1;97m"

    // High Intensity Backgrounds
    BgHiBlack   Color = "\033[0;100m"
    BgHiRed     Color = "\033[0;101m"
    BgHiGreen   Color = "\033[0;102m"
    BgHiYellow  Color = "\033[0;103m"
    BgHiBlue    Color = "\033[0;104m"
    BgHiPurple  Color = "\033[0;105m"
    BgHiCyan    Color = "\033[0;106m"
    BgHiWhite   Color = "\033[0;107m"

    // Styles
    Bold          Color = "\033[1m"
    Italic        Color = "\033[3m"
    Underline     Color = "\033[4m"
    Strikethrough Color = "\033[9m"
)


func (c Color) Apply(s string) string {
    return string(c) + s + string(Reset)
}
