// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Register(title string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div hx-ext=\"response-targets\" class=\"mt-[-12vh]\"><form hx-post=\"/register\" hx-trigger=\"submit\" hx-target-401=\"#register-error\" hx-swap=\"none\" hx-encoding=\"application/x-www-form-urlencoded\" class=\"flex flex-col w-[300px] bg-slate-100 p-4 pb-5 rounded-lg shadow-md gap-4\" id=\"form\"><div class=\"flex justify-center items-center border-b border-slate-300 pb-4\"><h1 class=\"text-xl font-bold text-slate-700\">Register</h1></div><div id=\"register-error\" class=\"text-red-500 text-sm\"></div><div class=\"flex flex-col gap-1\"><label for=\"Email\" class=\"text-slate-600 font-semibold\">Email</label> <input type=\"email\" name=\"Email\" id=\"email\" placeholder=\"name@company.com\" autocomplete=\"email\" class=\"rounded p-1 border border-gray-300\" required></div><div class=\"flex flex-col gap-1\"><label for=\"Username\" class=\"text-slate-600 font-semibold\">Username</label> <input type=\"text\" name=\"Username\" id=\"username\" placeholder=\"username\" class=\"rounded p-1 border border-gray-300\" required onkeyup=\"this.setCustomValidity(&#39;&#39;)\" hx-on:htmx:validation:validate=\"\r\n\t\t\t\t\t\tif (this.value.trim().length &lt; 4) {\r\n\t\t\t\t\t\t\tthis.setCustomValidity(&#39;Username must be at least 4 characters long.&#39;);\r\n\t\t\t\t\t\t\thtmx.find(&#39;#form&#39;).reportValidity();\r\n\t\t\t\t\t\t} else {\r\n\t\t\t\t\t\t\tthis.setCustomValidity(&#39;&#39;);\r\n\t\t\t\t\t\t\thtmx.find(&#39;#form&#39;).reportValidity();\r\n\t\t\t\t\t\t}\r\n\t\t\t\t\t\"></div><div class=\"flex flex-col gap-1\"><label for=\"Password\" class=\"text-slate-600 font-semibold\">Password</label> <input type=\"password\" name=\"Password\" id=\"password\" placeholder=\"••••••••\" autocomplete=\"new-password\" class=\"rounded p-1 border border-gray-300\" required onkeyup=\"this.setCustomValidity(&#39;&#39;)\" hx-on:htmx:validation:validate=\"\r\n\t\t\t\t\t\tif (this.value.length &lt; 8) {\r\n\t\t\t\t\t\t\tthis.setCustomValidity(&#39;Password must be at least 8 characters long.&#39;);\r\n\t\t\t\t\t\t\thtmx.find(&#39;#form&#39;).reportValidity();\r\n\t\t\t\t\t\t} else {\r\n\t\t\t\t\t\t\tthis.setCustomValidity(&#39;&#39;);\r\n\t\t\t\t\t\t\thtmx.find(&#39;#form&#39;).reportValidity();\r\n\t\t\t\t\t\t}\r\n\t\t\t\t\t\"></div><div class=\"flex flex-col gap-1\"><label for=\"ConfirmPassword\" class=\"text-slate-600 font-semibold\">Confirm Password</label> <input type=\"password\" name=\"ConfirmPassword\" id=\"confirm-password\" placeholder=\"••••••••\" autocomplete=\"new-password\" class=\"rounded p-1 border border-gray-300\" required onkeyup=\"this.setCustomValidity(&#39;&#39;)\" hx-on:htmx:validation:validate=\"\r\n\t\t\t\t\t\tif (this.value.length &lt; 8) {\r\n\t\t\t\t\t\t\tthis.setCustomValidity(&#39;Password must be at least 8 characters long.&#39;);\r\n\t\t\t\t\t\t\thtmx.find(&#39;#form&#39;).reportValidity();\r\n\t\t\t\t\t\t} else {\r\n\t\t\t\t\t\t\tthis.setCustomValidity(&#39;&#39;);\r\n\t\t\t\t\t\t\thtmx.find(&#39;#form&#39;).reportValidity();\r\n\t\t\t\t\t\t}\r\n\t\t\t\t\t\"></div><button type=\"submit\" class=\"text-white bg-amber-400 hover:bg-amber-500 font-semibold p-1 rounded-lg shadow-md transition duration-300\">Register</button><div class=\"flex justify-center text-center\"><p class=\"font-thin text-slate-600 text-sm\">Already have an account? <a href=\"/login\" class=\"text-amber-600 underline\">Login</a></p></div></form></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
