// This file contains the methods that panics when error return value is not nil.

package rod

import (
	"time"

	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// Connect to the browser and start to control it.
// If fails to connect, try to run a local browser, if local browser not found try to download one.
func (b *Browser) Connect() *Browser {
	utils.E(b.ConnectE())
	return b
}

// Close the browser and release related resources
func (b *Browser) Close() {
	_ = b.CloseE()
}

// Incognito creates a new incognito browser
func (b *Browser) Incognito() *Browser {
	b, err := b.IncognitoE()
	utils.E(err)
	return b
}

// Page creates a new tab
// If url is empty, the default target will be "about:blank".
func (b *Browser) Page(url string) *Page {
	p, err := b.PageE(url)
	utils.E(err)
	return p
}

// Pages returns all visible pages
func (b *Browser) Pages() Pages {
	list, err := b.PagesE()
	utils.E(err)
	return list
}

// PageFromTargetID creates a Page instance from a targetID
func (b *Browser) PageFromTargetID(targetID proto.TargetTargetID) *Page {
	p, err := b.PageFromTargetIDE(targetID)
	utils.E(err)
	return p
}

// HandleAuth for the next basic HTTP authentication.
// It will prevent the popup that requires user to input user name and password.
// Ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication
func (b *Browser) HandleAuth(username, password string) {
	wait := b.HandleAuthE(username, password)
	go func() { utils.E(wait()) }()
}

// FindByURL returns the page that has the url that matches the regex
func (ps Pages) FindByURL(regex string) *Page {
	p, err := ps.FindByURLE(regex)
	utils.E(err)
	return p
}

// Info of the page, such as the URL or title of the page
func (p *Page) Info() *proto.TargetTargetInfo {
	info, err := p.InfoE()
	utils.E(err)
	return info
}

// Cookies returns the page cookies. By default it will return the cookies for current page.
// The urls is the list of URLs for which applicable cookies will be fetched.
func (p *Page) Cookies(urls ...string) []*proto.NetworkCookie {
	cookies, err := p.CookiesE(urls)
	utils.E(err)
	return cookies
}

// SetCookies of the page.
// Cookie format: https://chromedevtools.github.io/devtools-protocol/tot/Network#method-setCookie
func (p *Page) SetCookies(cookies ...*proto.NetworkCookieParam) *Page {
	utils.E(p.SetCookiesE(cookies))
	return p
}

// SetExtraHeaders whether to always send extra HTTP headers with the requests from this page.
// The arguments are key-value pairs, you can set multiple key-value pairs at the same time.
func (p *Page) SetExtraHeaders(dict ...string) (cleanup func()) {
	cleanup, err := p.SetExtraHeadersE(dict)
	utils.E(err)
	return
}

// SetUserAgent Allows overriding user agent with the given string.
// If req is nil, the default user agent will be the same as a mac chrome.
func (p *Page) SetUserAgent(req *proto.NetworkSetUserAgentOverride) *Page {
	utils.E(p.SetUserAgentE(req))
	return p
}

// Navigate to url
// If url is empty, it will navigate to "about:blank".
func (p *Page) Navigate(url string) *Page {
	utils.E(p.NavigateE(url))
	return p
}

// GetWindow get window bounds
func (p *Page) GetWindow() *proto.BrowserBounds {
	bounds, err := p.GetWindowE()
	utils.E(err)
	return bounds
}

// Window set the window location and size
func (p *Page) Window(left, top, width, height int64) *Page {
	utils.E(p.WindowE(&proto.BrowserBounds{
		Left:        left,
		Top:         top,
		Width:       width,
		Height:      height,
		WindowState: proto.BrowserWindowStateNormal,
	}))
	return p
}

// WindowMinimize the window
func (p *Page) WindowMinimize() *Page {
	utils.E(p.WindowE(&proto.BrowserBounds{
		WindowState: proto.BrowserWindowStateMinimized,
	}))
	return p
}

// WindowMaximize the window
func (p *Page) WindowMaximize() *Page {
	utils.E(p.WindowE(&proto.BrowserBounds{
		WindowState: proto.BrowserWindowStateMaximized,
	}))
	return p
}

// WindowFullscreen the window
func (p *Page) WindowFullscreen() *Page {
	utils.E(p.WindowE(&proto.BrowserBounds{
		WindowState: proto.BrowserWindowStateFullscreen,
	}))
	return p
}

// WindowNormal the window size
func (p *Page) WindowNormal() *Page {
	utils.E(p.WindowE(&proto.BrowserBounds{
		WindowState: proto.BrowserWindowStateNormal,
	}))
	return p
}

// Viewport overrides the values of device screen dimensions.
func (p *Page) Viewport(width, height int64, deviceScaleFactor float64, mobile bool) *Page {
	utils.E(p.ViewportE(&proto.EmulationSetDeviceMetricsOverride{
		Width:             width,
		Height:            height,
		DeviceScaleFactor: deviceScaleFactor,
		Mobile:            mobile,
	}))
	return p
}

// Emulate the device, such as iPhone9. If device is empty, it will clear the override.
func (p *Page) Emulate(device devices.DeviceType) *Page {
	utils.E(p.EmulateE(device, false))
	return p
}

// StopLoading forces the page stop all navigations and pending resource fetches.
func (p *Page) StopLoading() *Page {
	utils.E(p.StopLoadingE())
	return p
}

// Close page
func (p *Page) Close() {
	utils.E(p.CloseE())
}

// HandleDialog accepts or dismisses next JavaScript initiated dialog (alert, confirm, prompt, or onbeforeunload)
// Because alert will block js, usually you have to run the wait function inside a goroutine. Check the unit test
// for it for more information.
func (p *Page) HandleDialog(accept bool, promptText string) (wait func()) {
	w := p.HandleDialogE(accept, promptText)
	return func() {
		utils.E(w())
	}
}

// Screenshot the page and returns the binary of the image
// If the toFile is "", it will save output to "tmp/screenshots" folder, time as the file name.
func (p *Page) Screenshot(toFile ...string) []byte {
	bin, err := p.ScreenshotE(false, &proto.PageCaptureScreenshot{})
	utils.E(err)
	utils.E(saveScreenshot(bin, toFile))
	return bin
}

// ScreenshotFullPage including all scrollable content and returns the binary of the image.
func (p *Page) ScreenshotFullPage(toFile ...string) []byte {
	bin, err := p.ScreenshotE(true, &proto.PageCaptureScreenshot{})
	utils.E(err)
	utils.E(saveScreenshot(bin, toFile))
	return bin
}

// PDF prints page as PDF
func (p *Page) PDF() []byte {
	pdf, err := p.PDFE(&proto.PagePrintToPDF{})
	utils.E(err)
	return pdf
}

// WaitOpen waits for a new page opened by the current one
func (p *Page) WaitOpen() (wait func() (newPage *Page)) {
	w := p.WaitOpenE()
	return func() *Page {
		page, err := w()
		utils.E(err)
		return page
	}
}

// WaitPauseOpen waits for a page opened by the current page, before opening pause the js execution.
// Because the js will be paused, you should put the code that triggers it in a goroutine, such as the click.
func (p *Page) WaitPauseOpen() (wait func() *Page, resume func()) {
	newPage, r, err := p.WaitPauseOpenE()
	utils.E(err)

	return func() *Page {
		page, err := newPage()
		utils.E(err)
		return page
	}, func() { utils.E(r()) }
}

// Pause stops on the next JavaScript statement
func (p *Page) Pause() *Page {
	utils.E(p.PauseE())
	return p
}

// WaitRequestIdle returns a wait function that waits until the page doesn't send request for 300ms.
// You can pass regular expressions to exclude the requests by their url.
func (p *Page) WaitRequestIdle(excludes ...string) (wait func()) {
	return p.WaitRequestIdleE(300*time.Millisecond, []string{""}, excludes)
}

// WaitIdle wait until the next window.requestIdleCallback is called.
func (p *Page) WaitIdle() *Page {
	utils.E(p.WaitIdleE(time.Minute))
	return p
}

// WaitLoad wait until the `window.onload` is complete, resolve immediately if already fired.
func (p *Page) WaitLoad() *Page {
	utils.E(p.WaitLoadE())
	return p
}

// AddScriptTag to page. If url is empty, content will be used.
func (p *Page) AddScriptTag(url string) *Page {
	utils.E(p.AddScriptTagE(url, ""))
	return p
}

// AddStyleTag to page. If url is empty, content will be used.
func (p *Page) AddStyleTag(url string) *Page {
	utils.E(p.AddStyleTagE(url, ""))
	return p
}

// EvalOnNewDocument Evaluates given script in every frame upon creation (before loading frame's scripts).
func (p *Page) EvalOnNewDocument(js string) {
	_, err := p.EvalOnNewDocumentE(js)
	utils.E(err)
}

// Expose function to the page's window object. Must bind before navigate to the page. Bindings survive reloads.
// Binding function takes exactly one argument, this argument should be string.
func (p *Page) Expose(name string) (callback chan string, stop func()) {
	c, s, err := p.ExposeE(name)
	utils.E(err)
	return c, s
}

// Eval js on the page. The first param must be a js function definition.
// For example page.Eval(`n => n + 1`, 1) will return 2
func (p *Page) Eval(js string, params ...interface{}) proto.JSON {
	res, err := p.EvalE(true, "", js, params)
	utils.E(err)
	return res.Value
}

// Wait js function until it returns true
func (p *Page) Wait(js string, params ...interface{}) {
	utils.E(p.WaitE(Sleeper(), "", js, params))
}

// ObjectToJSON by remote object
func (p *Page) ObjectToJSON(obj *proto.RuntimeRemoteObject) proto.JSON {
	j, err := p.ObjectToJSONE(obj)
	utils.E(err)
	return j
}

// ObjectsToJSON by remote objects
func (p *Page) ObjectsToJSON(list []*proto.RuntimeRemoteObject) proto.JSON {
	result := "[]"
	for _, obj := range list {
		j, err := p.ObjectToJSONE(obj)
		utils.E(err)
		result, err = sjson.SetRaw(result, "-1", j.Raw)
		utils.E(err)
	}
	return proto.JSON{Result: gjson.Parse(result)}
}

// ElementFromNode creates an Element from the node id
func (p *Page) ElementFromNode(id proto.DOMNodeID) *Element {
	el, err := p.ElementFromNodeE(id)
	utils.E(err)
	return el
}

// ElementFromPoint creates an Element from the absolute point on the page.
// The point should include the window scroll offset.
func (p *Page) ElementFromPoint(left, top int) *Element {
	el, err := p.ElementFromPointE(int64(left), int64(top))
	utils.E(err)
	return el
}

// Release remote object
func (p *Page) Release(objectID proto.RuntimeRemoteObjectID) *Page {
	utils.E(p.ReleaseE(objectID))
	return p
}

// Has an element that matches the css selector
func (p *Page) Has(selector string) bool {
	has, err := p.HasE(selector)
	utils.E(err)
	return has
}

// HasX an element that matches the XPath selector
func (p *Page) HasX(selector string) bool {
	has, err := p.HasXE(selector)
	utils.E(err)
	return has
}

// HasMatches an element that matches the css selector and its text matches the regex.
func (p *Page) HasMatches(selector, regex string) bool {
	has, err := p.HasMatchesE(selector, regex)
	utils.E(err)
	return has
}

// Search for each given query in the DOM tree until find one, before that it will keep retrying.
// The query can be plain text or css selector or xpath.
// It will search nested iframes and shadow doms too.
func (p *Page) Search(queries ...string) *Element {
	list, err := p.SearchE(Sleeper(), queries, 0, 1)
	utils.E(err)
	return list.First()
}

// Element retries until an element in the page that matches one of the CSS selectors
func (p *Page) Element(selectors ...string) *Element {
	el, err := p.ElementE(Sleeper(), "", selectors)
	utils.E(err)
	return el
}

// ElementMatches retries until an element in the page that matches one of the pairs.
// Each pairs is a css selector and a regex. A sample call will look like page.ElementMatches("div", "click me").
// The regex is the js regex, not golang's.
func (p *Page) ElementMatches(pairs ...string) *Element {
	el, err := p.ElementMatchesE(Sleeper(), "", pairs)
	utils.E(err)
	return el
}

// ElementByJS retries until returns the element from the return value of the js function
func (p *Page) ElementByJS(js string, params ...interface{}) *Element {
	el, err := p.ElementByJSE(Sleeper(), "", js, params)
	utils.E(err)
	return el
}

// Elements returns all elements that match the css selector
func (p *Page) Elements(selector string) Elements {
	list, err := p.ElementsE("", selector)
	utils.E(err)
	return list
}

// ElementsX returns all elements that match the XPath selector
func (p *Page) ElementsX(xpath string) Elements {
	list, err := p.ElementsXE("", xpath)
	utils.E(err)
	return list
}

// ElementX retries until an element in the page that matches one of the XPath selectors
func (p *Page) ElementX(xPaths ...string) *Element {
	el, err := p.ElementXE(Sleeper(), "", xPaths)
	utils.E(err)
	return el
}

// ElementsByJS returns the elements from the return value of the js
func (p *Page) ElementsByJS(js string, params ...interface{}) Elements {
	list, err := p.ElementsByJSE("", js, params)
	utils.E(err)
	return list
}

// Move to the absolute position
func (m *Mouse) Move(x, y float64) *Mouse {
	utils.E(m.MoveE(x, y, 0))
	return m
}

// Scroll with the relative offset
func (m *Mouse) Scroll(x, y float64) *Mouse {
	utils.E(m.ScrollE(x, y, 0))
	return m
}

// Down holds the button down
func (m *Mouse) Down(button proto.InputMouseButton) *Mouse {
	utils.E(m.DownE(button, 1))
	return m
}

// Up release the button
func (m *Mouse) Up(button proto.InputMouseButton) *Mouse {
	utils.E(m.UpE(button, 1))
	return m
}

// Click will press then release the button
func (m *Mouse) Click(button proto.InputMouseButton) *Mouse {
	utils.E(m.ClickE(button))
	return m
}

// Down holds key down
func (k *Keyboard) Down(key rune) *Keyboard {
	utils.E(k.DownE(key))
	return k
}

// Up releases the key
func (k *Keyboard) Up(key rune) *Keyboard {
	utils.E(k.UpE(key))
	return k
}

// Press a key
func (k *Keyboard) Press(key rune) *Keyboard {
	utils.E(k.PressE(key))
	return k
}

// InsertText like paste text into the page
func (k *Keyboard) InsertText(text string) *Keyboard {
	utils.E(k.InsertTextE(text))
	return k
}

// Describe returns the element info
// Returned json: https://chromedevtools.github.io/devtools-protocol/tot/DOM#type-Node
func (el *Element) Describe() *proto.DOMNode {
	node, err := el.DescribeE(1, false)
	utils.E(err)
	return node
}

// NodeID of the node
func (el *Element) NodeID() proto.DOMNodeID {
	id, err := el.NodeIDE()
	utils.E(err)
	return id
}

// ShadowRoot returns the shadow root of this element
func (el *Element) ShadowRoot() *Element {
	node, err := el.ShadowRootE()
	utils.E(err)
	return node
}

// Focus sets focus on the specified element
func (el *Element) Focus() *Element {
	utils.E(el.FocusE())
	return el
}

// ScrollIntoView scrolls the current element into the visible area of the browser
// window if it's not already within the visible area.
func (el *Element) ScrollIntoView() *Element {
	utils.E(el.ScrollIntoViewE())
	return el
}

// Hover the mouse over the center of the element.
func (el *Element) Hover() *Element {
	utils.E(el.HoverE())
	return el
}

// Click the element
func (el *Element) Click() *Element {
	utils.E(el.ClickE(proto.InputMouseButtonLeft))
	return el
}

// Clickable checks if the element is behind another element, such as when covered by a modal.
func (el *Element) Clickable() bool {
	clickable, err := el.ClickableE()
	utils.E(err)
	return clickable
}

// Press a key
func (el *Element) Press(key rune) *Element {
	utils.E(el.PressE(key))
	return el
}

// SelectText selects the text that matches the regular expression
func (el *Element) SelectText(regex string) *Element {
	utils.E(el.SelectTextE(regex))
	return el
}

// SelectAllText selects all text
func (el *Element) SelectAllText() *Element {
	utils.E(el.SelectAllTextE())
	return el
}

// Input wll click the element and input the text.
// To empty the input you can use something like el.SelectAllText().Input("")
func (el *Element) Input(text string) *Element {
	utils.E(el.InputE(text))
	return el
}

// Blur will call the blur function on the element.
// On inputs, this will deselect the element.
func (el *Element) Blur() *Element {
	utils.E(el.BlurE())
	return el
}

// Select the option elements that match the selectors, the selector can be text content or css selector
func (el *Element) Select(selectors ...string) *Element {
	utils.E(el.SelectE(selectors))
	return el
}

// Matches checks if the element can be selected by the css selector
func (el *Element) Matches(selector string) bool {
	res, err := el.MatchesE(selector)
	utils.E(err)
	return res
}

// Attribute returns the value of a specified attribute on the element.
// Please check the Property function before you use it, usually you don't want to use Attribute.
// https://stackoverflow.com/questions/6003819/what-is-the-difference-between-properties-and-attributes-in-html
func (el *Element) Attribute(name string) *string {
	attr, err := el.AttributeE(name)
	utils.E(err)
	return attr
}

// Property returns the value of a specified property on the element.
// It's similar to Attribute but attributes can only be string, properties can be types like bool, float, etc.
// https://stackoverflow.com/questions/6003819/what-is-the-difference-between-properties-and-attributes-in-html
func (el *Element) Property(name string) proto.JSON {
	prop, err := el.PropertyE(name)
	utils.E(err)
	return prop
}

// ContainsElement check if the target is equal or inside the element.
func (el *Element) ContainsElement(target *Element) bool {
	contains, err := el.ContainsElementE(target)
	utils.E(err)
	return contains
}

// SetFiles sets files for the given file input element
func (el *Element) SetFiles(paths ...string) *Element {
	utils.E(el.SetFilesE(paths))
	return el
}

// Text gets the innerText of the element
func (el *Element) Text() string {
	s, err := el.TextE()
	utils.E(err)
	return s
}

// HTML gets the outerHTML of the element
func (el *Element) HTML() string {
	s, err := el.HTMLE()
	utils.E(err)
	return s
}

// Visible returns true if the element is visible on the page
func (el *Element) Visible() bool {
	v, err := el.VisibleE()
	utils.E(err)
	return v
}

// WaitLoad for element like <img />
func (el *Element) WaitLoad() *Element {
	utils.E(el.WaitLoadE())
	return el
}

// WaitStable waits until the size and position are stable. Useful when waiting for the animation of modal
// or button to complete so that we can simulate the mouse to move to it and click on it.
func (el *Element) WaitStable() *Element {
	utils.E(el.WaitStableE(100 * time.Millisecond))
	return el
}

// Wait until the js returns true
func (el *Element) Wait(js string, params ...interface{}) *Element {
	utils.E(el.WaitE(js, params))
	return el
}

// WaitVisible until the element is visible
func (el *Element) WaitVisible() *Element {
	utils.E(el.WaitVisibleE())
	return el
}

// WaitInvisible until the element is not visible or removed
func (el *Element) WaitInvisible() *Element {
	utils.E(el.WaitInvisibleE())
	return el
}

// Box returns the size of an element and its position relative to the main frame.
func (el *Element) Box() *proto.DOMRect {
	box, err := el.BoxE()
	utils.E(err)
	return box
}

// CanvasToImage get image data of a canvas.
// The default format is image/png.
// The default quality is 0.92.
// doc: https://developer.mozilla.org/en-US/docs/Web/API/HTMLCanvasElement/toDataURL
func (el *Element) CanvasToImage(format string, quality float64) []byte {
	bin, err := el.CanvasToImageE(format, quality)
	utils.E(err)
	return bin
}

// Resource returns the binary of the "src" properly, such as the image or audio file.
func (el *Element) Resource() []byte {
	bin, err := el.ResourceE()
	utils.E(err)
	return bin
}

// Screenshot of the area of the element
func (el *Element) Screenshot(toFile ...string) []byte {
	bin, err := el.ScreenshotE(proto.PageCaptureScreenshotFormatPng, 0)
	utils.E(err)
	utils.E(saveScreenshot(bin, toFile))
	return bin
}

// Release remote object on browser
func (el *Element) Release() {
	utils.E(el.ReleaseE())
}

// Eval evaluates js function on the element, the first param must be a js function definition
// For example: el.Eval(`name => this.getAttribute(name)`, "value")
func (el *Element) Eval(js string, params ...interface{}) proto.JSON {
	res, err := el.EvalE(true, js, params)
	utils.E(err)
	return res.Value
}

// Has an element that matches the css selector
func (el *Element) Has(selector string) bool {
	has, err := el.HasE(selector)
	utils.E(err)
	return has
}

// HasX an element that matches the XPath selector
func (el *Element) HasX(selector string) bool {
	has, err := el.HasXE(selector)
	utils.E(err)
	return has
}

// HasMatches an element that matches the css selector and its text matches the regex.
func (el *Element) HasMatches(selector, regex string) bool {
	has, err := el.HasMatchesE(selector, regex)
	utils.E(err)
	return has
}

// Element returns the first child that matches the css selector
func (el *Element) Element(selector string) *Element {
	el, err := el.ElementE(selector)
	utils.E(err)
	return el
}

// ElementX returns the first child that matches the XPath selector
func (el *Element) ElementX(xpath string) *Element {
	el, err := el.ElementXE(xpath)
	utils.E(err)
	return el
}

// ElementByJS returns the element from the return value of the js
func (el *Element) ElementByJS(js string, params ...interface{}) *Element {
	el, err := el.ElementByJSE(js, params)
	utils.E(err)
	return el
}

// Parent returns the parent element
func (el *Element) Parent() *Element {
	parent, err := el.ParentE()
	utils.E(err)
	return parent
}

// Parents that match the selector
func (el *Element) Parents(selector string) Elements {
	list, err := el.ParentsE(selector)
	utils.E(err)
	return list
}

// Next returns the next sibling element
func (el *Element) Next() *Element {
	parent, err := el.NextE()
	utils.E(err)
	return parent
}

// Previous returns the previous sibling element
func (el *Element) Previous() *Element {
	parent, err := el.PreviousE()
	utils.E(err)
	return parent
}

// ElementMatches returns the first element in the page that matches the CSS selector and its text matches the regex.
// The regex is the js regex, not golang's.
func (el *Element) ElementMatches(selector, regex string) *Element {
	el, err := el.ElementMatchesE(selector, regex)
	utils.E(err)
	return el
}

// Elements returns all elements that match the css selector
func (el *Element) Elements(selector string) Elements {
	list, err := el.ElementsE(selector)
	utils.E(err)
	return list
}

// ElementsX returns all elements that match the XPath selector
func (el *Element) ElementsX(xpath string) Elements {
	list, err := el.ElementsXE(xpath)
	utils.E(err)
	return list
}

// ElementsByJS returns the elements from the return value of the js
func (el *Element) ElementsByJS(js string, params ...interface{}) Elements {
	list, err := el.ElementsByJSE(js, params)
	utils.E(err)
	return list
}
