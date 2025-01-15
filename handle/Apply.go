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
	task any
}

type TaskUpsertZone struct {
	zone *cli.Zone
}

type TaskSetCollectorDisplay struct {
	collectorKey string
	display      bool
}

type TaskCreateAsset struct {
	asset *cli.AssetCli
}

type TaskSetAssetMode struct {
	asset *cli.AssetCli
}

type TaskSetAssetKind struct {
	asset *cli.AssetCli
}

type TaskSetAssetZone struct {
	asset *cli.AssetCli
}

type TaskDeleteLabelFromAsset struct {
	asset   *cli.AssetCli
	labelId int
}

type TaskEnableAssetCheck struct {
	asset        *cli.AssetCli
	collectorKey string
	checkKey     string
}

func processChanges(api, token string, containerId int, changes *[]*Change) {
	n := len(*changes)
	for i, c := range *changes {
		fmt.Printf("Processing task %d/%d: %s ...\n", i+1, n, c.info)
		var err error
		switch task := c.task.(type) {
		case TaskUpsertZone:
			err = req.UpsertZone(api, token, containerId, task.zone.Zone, task.zone.Name)
		case TaskSetCollectorDisplay:
			err = req.SetCollectorDisplay(api, token, containerId, task.collectorKey, task.display)
		case TaskCreateAsset:
			task.asset.Id, err = req.CreateAsset(api, token, containerId, task.asset.Name)
		case TaskSetAssetMode:
			err = req.SetAssetMode(api, token, task.asset.Id, task.asset.Mode, nil)
		case TaskSetAssetKind:
			err = req.SetAssetKind(api, token, task.asset.Id, task.asset.Kind)
		case TaskSetAssetZone:
			err = req.SetAssetZone(api, token, task.asset.Id, *task.asset.Zone)
		case TaskDeleteLabelFromAsset:
			err = req.DeleteLabelFromAsset(api, token, task.asset.Id, task.labelId)
		case TaskEnableAssetCheck:
			err = req.EnableAssetCheck(api, token, task.asset.Id, task.collectorKey, task.checkKey)
		}
		util.ExitOnErr(err)
	}
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
			// TODO: Rename label API
		})
	}
}

func assetChanges(changes *[]*Change, purge bool, ca, ta *cli.AssetCli, cs, ts *cli.State) {
	if ta.Mode != "" && ta.Mode != ca.Mode {
		*changes = append(*changes, &Change{
			info: fmt.Sprintf("Set asset mode for asset '%s' to: %s", cval(ta.Str()), cval(ta.Mode)),
			task: TaskSetAssetMode{asset: ta},
		})
	}
	if ta.Kind != "" && ta.Kind != ca.Kind {
		*changes = append(*changes, &Change{
			info: fmt.Sprintf("Set asset kind for asset '%s' to: %s", cval(ta.Str()), cval(ta.Kind)),
			task: TaskSetAssetKind{asset: ta},
		})
	}
	if ta.Zone != nil && (ca.Zone == nil || *ta.Zone != *ca.Zone) {
		*changes = append(*changes, &Change{
			info: fmt.Sprintf("Set asset zone for asset '%s' to: %s", cval(ta.Str()), cval(*ta.Zone)),
			task: TaskSetAssetZone{asset: ta},
		})
	}
	if purge {
		if ca.Labels != nil && ta.Labels != nil {
			for _, key := range *ca.Labels {
				if label := cs.LabelByKey(key); label != nil {
					if !ta.HasLabelId(label.Id, ts.GetLabelMap()) {
						*changes = append(*changes, &Change{
							info: fmt.Sprintf("Delete label '%s' from asset '%s'", cval(label.Str()), cval(ta.Str())),
							task: TaskDeleteLabelFromAsset{asset: ta, labelId: label.Id},
						})
					}
				}
			}
		}
		if ca.DisabledChecks != nil && ta.DisabledChecks != nil {
			for _, disabledChk := range *ca.DisabledChecks {
				found := false
				for _, c := range *ta.DisabledChecks {
					if c.Collector == disabledChk.Collector && c.Check == disabledChk.Check {
						found = true
						break
					}
				}
				if !found {
					*changes = append(*changes, &Change{
						info: fmt.Sprintf("Enable collector check '%s/%s' on asset '%s'", cval(disabledChk.Collector), cval(disabledChk.Check), cval(ta.Str())),
						task: TaskEnableAssetCheck{asset: ta, collectorKey: disabledChk.Collector, checkKey: disabledChk.Check},
					})
				}
			}
		}
	}
}

func ensureChanges(api, token string, purge bool, cs, ts *cli.State, cMap map[string]*cli.Collector) []*Change {
	changes := []*Change{}
	//
	// Container changes
	//
	if ts.Container.Name != "" && ts.Container.Name != cs.Container.Name {
		changes = append(changes, &Change{
			info: fmt.Sprintf("Rename container ID %s from '%s' to '%s'", cval(cs.Container.Id), cval(cs.Container.Name), cval(ts.Container.Name)),
			// TODO: Rename container API
		})
	}
	//
	// Show collectors and sanity check mode and disabled checks
	//
	enableCollector := cli.StrSet{}
	for _, ta := range ts.Assets {
		if ta.Collectors != nil {
			for _, c := range *ta.Collectors {
				if enableCollector.Has(c.Key) {
					continue
				}
				if _, ok := cMap[c.Key]; !ok {
					changes = append(changes, &Change{
						info: fmt.Sprintf("Enable collector: %s", cval(c.Key)),
						task: TaskSetCollectorDisplay{collectorKey: c.Key, display: true},
					})
					enableCollector.Set(c.Key)
				}
			}
		}
		if ta.DisabledChecks != nil {
			for _, disabledChk := range *ta.DisabledChecks {
				found := false
				if ta.Collectors != nil {
					for _, c := range *ta.Collectors {
						if c.Key == disabledChk.Collector {
							found = true
							break
						}
					}
					if !found {
						util.ExitErr("Collector '%s' is not configured for asset '%s', but a disabled check for it exists.", disabledChk.Collector, ta.Str())
					}
				}

			}
		}
		switch ta.Mode {
		case "", "normal", "maintenance", "disabled":
			continue
		}
		util.ExitErr("Asset '%s' has an invalid mode '%s'. Must be one of {normal,maintenance,disabled}", ta.Str(), ta.Mode)
	}

	//
	// Zone changes
	//
	for _, tz := range ts.Zones {
		cz := cs.ZoneById(tz.Zone)
		if cz == nil {
			if tz.Zone < 1 || tz.Zone > 9 {
				util.ExitErr("Invalid zone '%d'. Must be a value between 1 and 9.", tz.Zone)
			}
			if tz.Name == "" {
				util.ExitErr("Zone '%d' is new and therefore requires a name", tz.Zone)
			}
			changes = append(changes, &Change{
				info: fmt.Sprintf("Create new zone: %s", cval(tz.Str())),
				task: TaskUpsertZone{zone: tz},
			})
		} else {
			if tz.Name != "" && tz.Name != cz.Name {
				changes = append(changes, &Change{
					info: fmt.Sprintf("Rename zone ID %s from '%s' to '%s'", cval(tz.Zone), cval(cz.Name), cval(tz.Name)),
					task: TaskUpsertZone{zone: tz},
				})
			}
		}
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
				// TODO: Create new label API
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
	for _, ta := range ts.Assets {
		if ta.Labels != nil {
			for _, labelKey := range *ta.Labels {
				if _, ok := ts.Labels[labelKey]; !ok {
					util.ExitErr("Asset '%s' is using label reference '%s' which does not exist in 'labels'.", ta.Str(), labelKey)
				}
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
				task: TaskCreateAsset{asset: ta},
			})
			assetChanges(&changes, purge, &cli.DefaultAsset, ta, cs, ts)
		} else {
			ca := cs.AssetById(ta.Id)
			if ca == nil {
				util.ExitErr("Asset ID %d not found in container '%s'.", ta.Id, cs.Container.Str())
			}
			assetChanges(&changes, purge, ca, ta, cs, ts)
		}
	}
	return changes
}

func getCacheState(containerId int) *cli.State {
	state := cli.StateFromCache(containerId)
	if state != nil {
		if age, err := state.GetAge(); err == nil {
			if age.Hours() < 8 {
				util.Color("A cache for container ID %d was found that is only %s old. Would you like to use it? (yes/no): ", containerId, util.HumanizeDuration(*age))

				fmt.Println()
				return state
				// TODO ask ->
				// if util.AskForConfirmation() {
				// 	return state
				// }
			}
		}
	}
	return nil
}

func sanityCheckCollectorConfig(ts *cli.State, cMap map[string]*cli.Collector) {
	for _, asset := range ts.Assets {
		if asset.Collectors == nil {
			continue
		}
		for _, collector := range *asset.Collectors {
			if c, ok := cMap[collector.Key]; ok {
				for k, v := range collector.Config {
					found := false
					for _, o := range c.Options {
						if o.Key == k {
							found = true
							if o.Type == "String" && (k == "password" || k == "secret") {
								if _, ok := v.(string); ok {
									break
								}
								if obj, ok := v.(map[string]any); ok {
									if v, ok := obj["encrypted"]; ok && len(obj) == 1 {
										if _, ok := v.(string); ok {
											break
										}
									}
								}
								util.ExitErr("Collector '%s' on asset '%s' expects property '%s' to be a string or encryption value", collector.Key, asset.Str(), k)
							}
							switch o.Type {
							case "Bool":
								if _, ok := v.(bool); !ok {
									util.ExitErr("Collector '%s' on asset '%s' expects a boolean value for property '%s' but found type %T", collector.Key, asset.Str(), k, v)
								}
							case "Int":
								if _, ok := v.(int); !ok {
									util.ExitErr("Collector '%s' on asset '%s' expects an integer value for property '%s' but found type %T", collector.Key, asset.Str(), k, v)
								}
							case "Float":
								if _, ok := v.(float64); !ok {
									util.ExitErr("Collector '%s' on asset '%s' expects a floating point for property '%s' but found type %T", collector.Key, asset.Str(), k, v)
								}
							case "String":
								if _, ok := v.(string); !ok {
									util.ExitErr("Collector '%s' on asset '%s' expects a string value for property '%s' but found type %T", collector.Key, asset.Str(), k, v)
								}
							case "ListBool", "ListInt", "ListFloat", "ListString":
								if arr, ok := v.([]interface{}); ok {
									switch o.Type {
									case "ListBool":
										for _, v := range arr {
											if _, ok = v.(bool); !ok {
												util.ExitErr("Collector '%s' on asset '%s' expects a list of boolean values for property '%s' but the list contains type %T", collector.Key, asset.Str(), k, v)
											}
										}
									case "ListInt":
										for _, v := range arr {
											if _, ok = v.(int); !ok {
												util.ExitErr("Collector '%s' on asset '%s' expects a list of integer values for property '%s' but the list contains type %T", collector.Key, asset.Str(), k, v)
											}
										}
									case "ListFloat":
										for _, v := range arr {
											if _, ok = v.(float64); !ok {
												util.ExitErr("Collector '%s' on asset '%s' expects a list of floating point values for property '%s' but the list contains type %T", collector.Key, asset.Str(), k, v)
											}
										}
									case "ListString":
										for _, v := range arr {
											if _, ok = v.(string); !ok {
												util.ExitErr("Collector '%s' on asset '%s' expects a list of string values for property '%s' but the list contains type %T", collector.Key, asset.Str(), k, v)
											}
										}
									}
									break
								}
								util.ExitErr("Collector '%s' on asset '%s' expects a list of values for property '%s' but found type %T", collector.Key, asset.Str(), k, v)
							}
							break
						}
					}
					if !found {
						util.ExitErr("Collector '%s' on asset '%s' contains an unknown configuration property '%s'.", collector.Key, asset.Str(), k)
					}
				}
			}
		}
	}
}

func revertUse(assets []*cli.AssetCli) {
	for _, asset := range assets {
		if asset.Collectors != nil {
			for _, c := range *asset.Collectors {
				if v, ok := c.Config["use"]; ok {
					delete(c.Config, "use")
					c.Config["_use"] = v
				}
			}
		}
	}
}

func niceAssetKinds(api string, assets []*cli.AssetCli) {
	kinds, err := req.GetAssetKinds(api)
	util.ExitOnErr(err)

	for _, asset := range assets {
		if asset.Kind != "" {
			kind := util.InSlice(kinds, asset.Kind)
			if kind == nil {
				util.ExitErr("Asset '%s' has an invalid asset kind: %s", asset.Str(), asset.Kind)
			}
			asset.Kind = *kind
		}
	}
}

func Apply(api, token, filename string, dryRun, purge bool) {
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

	cs := getCacheState(ts.Container.Id)
	if cs == nil {
		fmt.Println("Read current state...")
		cs = ensureState(&TGetAssets{
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
		cs.WriteCache()
	}

	cMap := map[string]*cli.Collector{}

	if ts.HasAssetKind() {
		niceAssetKinds(api, ts.Assets)
	}

	if ts.HasCollector() {
		// use -> _use
		revertUse(cs.Assets)
		revertUse(ts.Assets)

		fmt.Println("Read collectors...")
		collectors, err := req.GetCollectors(api, token, ts.Container.Id, []string{"key"}, true)
		util.ExitOnErr(err)

		for _, c := range collectors {
			cMap[c.Key] = c
		}
		sanityCheckCollectorConfig(ts, cMap)
	}

	changes := ensureChanges(api, token, purge, cs, ts, cMap)
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
	} else {
		util.Color("Do you want to apply the change%s? (yes/no): ", util.Plural(n))
		if util.AskForConfirmation() {
			fmt.Println("")
			processChanges(api, token, ts.Container.Id, &changes)
			fmt.Println("")
			util.ExitOk("Done.")
		}
		util.ExitOk("Cancelled.")
	}
}
