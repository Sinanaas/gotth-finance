// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/Sinanaas/gotth-financial-tracker/internal/models"
import "github.com/Sinanaas/gotth-financial-tracker/internal/utils"
import "github.com/Sinanaas/gotth-financial-tracker/internal/constants"
import "fmt"
import "strconv"

func Loans(categories []models.Category, loans []models.Loan, accounts []models.Account, transactionType []string) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<main class=\"container mx-auto p-4 min-h-full top-0\"><div class=\"flex flex-col md:flex-row gap-6\"><div hx-ext=\"response-targets\" class=\"w-full md:w-1/3\"><form hx-post=\"/loan\" hx-trigger=\"submit\" hx-encoding=\"application/x-www-form-urlencoded\" class=\"flex flex-col bg-slate-100 p-6 rounded-lg shadow-md gap-4\"><div class=\"flex justify-center items-center border-b border-slate-300 pb-4 mb-4\"><h2 class=\"text-xl font-bold text-slate-700\">Add a Loan</h2></div><div class=\"flex flex-col gap-2\"><label for=\"amount\" class=\"text-slate-600 font-semibold\">Amount</label> <input type=\"number\" name=\"Amount\" id=\"amount\" placeholder=\"0.00\" class=\"rounded p-2 border border-slate-300\" min=\"0\" required onkeyup=\"this.setCustomValidity(&#39;&#39;)\" hx-on:htmx:validation:validate=\"\r\n                            if (this.value &lt; 0) {\r\n                                this.setCustomValidity(&#39;Amount cannot be negative.&#39;);\r\n                                htmx.find(&#39;#account-form&#39;).reportValidity();\r\n                            } else {\r\n                                this.setCustomValidity(&#39;&#39;);\r\n                                htmx.find(&#39;#account-form&#39;).reportValidity();\r\n                            }\r\n                        \"></div><div class=\"flex flex-col gap-2\"><label for=\"type\" class=\"text-slate-600 font-semibold\">Type</label> <select name=\"Type\" id=\"type\" class=\"rounded p-2 border border-slate-300\" required oninvalid=\"this.setCustomValidity(&#39;Please select a transaction type.&#39;)\" oninput=\"this.setCustomValidity(&#39;&#39;)\"><option disabled selected value class=\"w-full text-slate-600\">select a transaction type</option> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for i := 0; i < len(transactionType); i++ {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<option class=\"w-full\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var2 string
			templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.Itoa(i))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/loan.templ`, Line: 58, Col: 70}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if transactionType[i] == "Expenses" {
				var templ_7745c5c3_Var3 string
				templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(transactionType[i] + " / Memberi Hutang")
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/loan.templ`, Line: 60, Col: 70}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				var templ_7745c5c3_Var4 string
				templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(transactionType[i] + " / Hutang")
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/loan.templ`, Line: 62, Col: 62}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</option>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</select></div><div class=\"flex w-full justify-between gap-6\"><div class=\"flex flex-col gap-2 w-1/2\"><label for=\"account\" class=\"text-slate-600 font-semibold\">Account</label> <select name=\"Account\" id=\"account\" class=\"rounded p-2 w-full border border-slate-300\" hx-get=\"/account/balance\" hx-target=\"#balance\" hx-params=\"Account\" required oninvalid=\"this.setCustomValidity(&#39;Please select an account.&#39;)\" oninput=\"this.setCustomValidity(&#39;&#39;)\"><option disabled selected value class=\"w-full text-slate-600\">select an account</option> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, account := range accounts {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<option class=\"w-full\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%s", account.ID))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/loan.templ`, Line: 84, Col: 88}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var6 string
			templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(account.Name)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/loan.templ`, Line: 85, Col: 42}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</option>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</select></div><div class=\"flex flex-col gap-2\"><label for=\"balance\" class=\"text-slate-600 font-semibold\">Balance</label><div id=\"balance\" class=\"text-sm text-gray-900\"><div class=\"text-sm text-gray-900\"><input type=\"text\" value=\"\" disabled class=\"rounded p-2 border border-slate-300\"></div></div></div></div><div class=\"flex flex-col gap-2\"><label for=\"towhom\" class=\"text-slate-600 font-semibold\">To/From Whom</label> <input type=\"text\" name=\"Towhom\" id=\"towhom\" placeholder=\"To/From Whom?\" class=\"rounded p-2 border border-slate-300\" required oninvalid=\"this.setCustomValidity(&#39;Please enter a name.&#39;)\" oninput=\"this.setCustomValidity(&#39;&#39;)\"></div><div class=\"flex flex-col gap-2\"><label for=\"description\" class=\"text-slate-600 font-semibold\">Description</label> <input type=\"text\" name=\"Description\" id=\"description\" placeholder=\"Description\" class=\"rounded p-2 border border-slate-300\" required oninvalid=\"this.setCustomValidity(&#39;Please enter a description.&#39;)\" oninput=\"this.setCustomValidity(&#39;&#39;)\"></div><div class=\"flex flex-col gap-2\"><label for=\"date\" class=\"text-slate-600 font-semibold\">Loan Date</label> <input type=\"date\" name=\"Date\" id=\"date\" class=\"rounded p-2 border border-slate-300\" required oninvalid=\"this.setCustomValidity(&#39;Please select a date.&#39;)\" oninput=\"this.setCustomValidity(&#39;&#39;)\"></div><div class=\"flex flex-col gap-2\"><label for=\"category\" class=\"text-slate-600 font-semibold\">Category</label> <select name=\"Category\" id=\"category\" class=\"rounded p-2 border border-slate-300\" required oninvalid=\"this.setCustomValidity(&#39;Please select a category.&#39;)\" oninput=\"this.setCustomValidity(&#39;&#39;)\"><option disabled selected value class=\"w-full text-slate-600\">select a category</option> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, category := range categories {
			if category.Name != "Income" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<option value=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var7 string
				templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%s", category.ID))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/loan.templ`, Line: 155, Col: 70}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var8 string
				templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(category.Name)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/loan.templ`, Line: 155, Col: 88}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</option>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</select></div><button type=\"submit\" class=\"text-white bg-amber-400 hover:bg-amber-500 font-semibold p-2 rounded-lg shadow-md transition duration-300\">Add Transaction</button></form></div><div class=\"w-full md:w-2/3\"><div class=\"w-full flex justify-center items-center\"><h1 class=\"text-4xl text-white font-bold text-center mb-6\">HUTANG</h1></div><table class=\"min-w-full divide-y divide-gray-200\"><thead class=\"bg-gray-50\"><tr><th scope=\"col\" class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider\">Loan TANGAL</th><th scope=\"col\" class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider\">Amount</th><th scope=\"col\" class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider\">Description</th><th scope=\"col\" class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider\">To Whom</th><th scope=\"col\" class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider\">Type</th><th scope=\"col\" class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider\">Category</th><th scope=\"col\" class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider\">Account</th><th scope=\"col\" class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider text-center\">Actions</th></tr></thead> <tbody class=\"bg-white divide-y divide-gray-200\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, loan := range loans {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<tr><td class=\"px-6 py-4 whitespace-nowrap\"><div class=\"text-sm text-gray-900\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var9 string
			templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(loan.LoanDate.Format("2006-01-02"))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/loan.templ`, Line: 211, Col: 95}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></td>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if loan.TransactionType == constants.Income {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<td class=\"px-6 py-4 whitespace-nowrap\"><div class=\"text-sm text-green-200\">Rp. ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var10 string
				templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(utils.FormatCurrency(loan.Amount))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/loan.templ`, Line: 215, Col: 99}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></td>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<td class=\"px-6 py-4 whitespace-nowrap\"><div class=\"text-sm text-red-200\">Rp. ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var11 string
				templ_7745c5c3_Var11, templ_7745c5c3_Err = templ.JoinStringErrs(utils.FormatCurrency(loan.Amount))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/loan.templ`, Line: 219, Col: 97}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var11))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></td>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<td class=\"px-6 py-4 whitespace-nowrap\"><div class=\"text-sm text-gray-900\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var12 string
			templ_7745c5c3_Var12, templ_7745c5c3_Err = templ.JoinStringErrs(loan.Description)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/loan.templ`, Line: 223, Col: 77}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var12))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></td><td class=\"px-6 py-4 whitespace-nowrap\"><div class=\"text-sm text-gray-900\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var13 string
			templ_7745c5c3_Var13, templ_7745c5c3_Err = templ.JoinStringErrs(loan.ToWhom)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/loan.templ`, Line: 226, Col: 72}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var13))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></td><td class=\"px-6 py-4 whitespace-nowrap\"><div class=\"text-sm text-gray-900\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var14 string
			templ_7745c5c3_Var14, templ_7745c5c3_Err = templ.JoinStringErrs(loan.TransactionType.ToString())
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/loan.templ`, Line: 229, Col: 92}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var14))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></td><td class=\"px-6 py-4 whitespace-nowrap\"><div class=\"text-sm text-gray-900\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var15 string
			templ_7745c5c3_Var15, templ_7745c5c3_Err = templ.JoinStringErrs(loan.Category.Name)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/loan.templ`, Line: 232, Col: 79}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var15))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></td><td class=\"px-6 py-4 whitespace-nowrap\"><div class=\"text-sm text-gray-900\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var16 string
			templ_7745c5c3_Var16, templ_7745c5c3_Err = templ.JoinStringErrs(loan.Account.Name)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/loan.templ`, Line: 235, Col: 78}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var16))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></td><td class=\"px-6 py-4 whitespace-nowrap flex gap-2\"><form hx-post=\"/loan/finish\" hx-trigger=\"submit\" hx-encoding=\"application/x-www-form-urlencoded\" hx-params=\"LoanID\" class=\"flex justify-center mb-0\"><button type=\"submit\" class=\"text-white bg-amber-400 px-2 hover:bg-amber-500 font-semibold p-1 rounded-lg shadow-md transition duration-300\" name=\"LoanID\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var17 string
			templ_7745c5c3_Var17, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%s", loan.ID))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/loan.templ`, Line: 249, Col: 70}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var17))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">Finish</button></form><form hx-put=\"/loan\" hx-trigger=\"submit\" hx-target=\"closest tr\" hx-encoding=\"application/x-www-form-urlencoded\" class=\"flex justify-center mb-0\"><input type=\"hidden\" name=\"LoanID\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var18 string
			templ_7745c5c3_Var18, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%s", loan.ID))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/loan.templ`, Line: 261, Col: 97}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var18))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> <button type=\"submit\" class=\"text-white bg-red-400 px-2 hover:bg-red-500 font-semibold p-1 rounded-lg shadow-md transition duration-300\">Delete</button></form></td></tr>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</tbody></table></div></div></main>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
