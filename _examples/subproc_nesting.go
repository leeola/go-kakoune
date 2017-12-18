package _examples

// // SubprocNesting enables GoKakoune to dynamically specify the indexes of subprocs.
// //
// // Eg, calling the example binary with:
// //
// var SubprocNesting = []api.Subproc{
// 	api.Sh{
// 		Func: func(kak *api.Kak) error {
// 			kak.Echo("hello from 2nd func")
// 			return nil
// 		},
// 	},
// 	api.Prompt{
// 	api.Prompt("foo", api.Subproc{
//   	Prompt: "foo",
// 		Func: func(kak *api.Kak) error {
// 			kak.Echo("hello from 2nd func")
// 			return nil
// 		},
// 	}),
// }
//
// var SubprocNesting = []api.Subproc{
// 	{
// 		Func: func(kak *api.Kak) error {
// 			kak.Echo("hello from 2nd func")
// 			return nil
// 		},
// 	},
// 	{Text: "name:",
//
// 	{
// 		Func: func(kak *api.Kak) error {
// 			return kak.Prompt("foo", api.Subproc{
// 				Func: func(kak *api.Kak) error {
// 					kak.Echo("hello from 2nd func")
// 					return nil
// 				},
// 			})
// 		},
// 	},
// }
