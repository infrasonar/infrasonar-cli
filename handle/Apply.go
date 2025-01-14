package handle

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/infrasonar/infrasonar-cli/cli"
	"github.com/infrasonar/infrasonar-cli/handle/util"
	"github.com/infrasonar/infrasonar-cli/req"
)

type Change struct {
	info string
}

func cval(a any) string {
	return color.HiCyanString("%v", a)
}

func readLabelChanges(changes *[]*Change, cl, tl *cli.Label) {
	if cl.Id != tl.Id {
		panic("label ID mismatch")
	}
	if tl.Name != "" && tl.Name != cl.Name {
		*changes = append(*changes, &Change{
			info: fmt.Sprintf("Rename label ID %s from '%s' to '%s'", cval(cl.Id), cval(cl.Name), cval(tl.Name)),
		})
	}
}

func ensureChanges(api, token string, noRemove bool, cs, ts *cli.State) []*Change {
	changes := []*Change{}
	if ts.Container.Name != "" && ts.Container.Name != cs.Container.Name {
		changes = append(changes, &Change{
			info: fmt.Sprintf("Rename container ID %s from '%s' to '%s'", cval(cs.Container.Id), cval(cs.Container.Name), cval(ts.Container.Name)),
		})
	}
	missingLabelIds := cli.IntSet{}
	for _, tl := range ts.Labels {
		if tl.Id == 0 {
			// New label
			if tl.Name == "" {
				util.ExitErr("One or more labels are missing both an 'id' and a 'name'. At least one of these attributes is required for each label.")
			}
			changes = append(changes, &Change{
				info: fmt.Sprintf("Create new label: %s", cval(tl.Name)),
			})
		} else {
			// Label with ID
			if cl := cs.LabelById(tl.Id); cl == nil {
				missingLabelIds.Set(tl.Id)
			} else {
				readLabelChanges(&changes, cl, tl)
			}
		}
	}
	if len(missingLabelIds) > 0 {
		fmt.Println("Get missing labels...")
		lm, err := req.GetLabels(api, token, missingLabelIds)
		util.ExitOnErr(err)
		for _, tl := range ts.Labels {
			if tl.Id != 0 {
				if cl := cs.LabelById(tl.Id); cl == nil {
					if cl := lm.LabelById(tl.Id); cl != nil {
						readLabelChanges(&changes, cl, tl)
					}
				}
			}
		}
	}
	for _, tz := range ts.Zones {
		cz := cs.ZoneById(tz.Zone)
		if cz == nil {
			if tz.Zone < 1 || tz.Zone > 9 {
				util.ExitErr("Invalid zone '%d'. Must be a value between 1 and 9.", tz.Zone)
			}
			changes = append(changes, &Change{
				info: fmt.Sprintf("Create new zone: %s", cval(tz.Str())),
			})
		} else {
			if tz.Name != "" && tz.Name != cz.Name {
				changes = append(changes, &Change{
					info: fmt.Sprintf("Rename zone ID %s from '%s' to '%s'", cval(tz.Zone), cval(cz.Name), cval(tz.Name)),
				})
			}
		}
	}
	for _, ta := range ts.Assets {
		for _, labelKey := range ta.Labels {
			if _, ok := ts.Labels[labelKey]; !ok {
				util.ExitErr("Asset '%s' is using label reference '%s' which does not exist in 'labels'.", ta.Str(), labelKey)
			}
		}
		if ta.Zone != nil {
			if zone := ts.ZoneById(*ta.Zone); zone == nil {
				util.ExitErr("Asset '%s' is using zone ID %d which does not exist in 'zones'.", ta.Str(), *ta.Zone)
			}
		}

		if ta.Id == 0 {
			// New asset
			if ta.Name == "" {
				util.ExitErr("One or more assets are missing both an 'id' and a 'name'. At least one of these attributes is required for each asset.")
			}
			changes = append(changes, &Change{
				info: fmt.Sprintf("Create new asset: %s", cval(ta.Name)),
			})
		} else {
			ca := cs.AssetById(ta.Id)
			if ca == nil {
				util.ExitErr("Asset ID %d not found in container '%s'.", ta.Id, cs.Container.Str())
			}
		}
	}
	return changes
}

func Apply(api, token, filename string, dryRun, noRemove bool) {
	if dryRun {
		util.Color(`-----------------------------------------
  Simulation :: no changes will be made
-----------------------------------------
`)
	}
	fmt.Println("Read input file...")
	ts, err := cli.StateFromFile(filename)
	util.ExitOnErr(err)

	if ts.Container.Id == 0 {
		util.ExitErr("missing container ID in input file")
	}

	fmt.Println("Read current state...")
	cs := ensureState(&TGetAssets{
		Api:             api,
		Token:           token,
		Output:          "",
		OutFn:           "-", // Force progress output, nothing will be written
		Container:       ts.Container.Id,
		Asset:           0,
		Properties:      cli.AssetProperties,
		Filters:         []string{},
		IncludeDefaults: true,
	})

	changes := ensureChanges(api, token, noRemove, cs, ts)
	n := len(changes)

	if n == 0 {
		util.ExitOk("No changes found.")
	}

	util.Color("Found %d change%s. Show details? (yes/no): ", n, util.Plural(n))
	if util.AskForConfirmation() {
		fmt.Println("")
		for _, c := range changes {
			fmt.Printf("- %s\n", c.info)
		}
		fmt.Println("")
	}

	if dryRun {
		util.ExitOk("Done. (no changes made)")
	}
	util.ExitOk("Done.")
}
