import { useState } from "react";
import { StyleSheet, Switch, Text, View } from "react-native";
import { getMobileMemberships, getMobileResidences } from "@/features/auth/mobile-access";
import { useGetV1BootstrapQuery } from "@/lib/api/generated-api";
import { useAppSelector } from "@/redux/hooks";
import {
  type NotificationPushPreferences,
  useGetNotificationPushPreferencesQuery,
  usePutNotificationPushPreferencesMutation,
} from "@/lib/api/notification-api-extensions";

const options: {key:keyof NotificationPushPreferences;label:string}[] = [
  {key:"visitor_updates_push",label:"Visitor updates"},
  {key:"maintenance_push",label:"Maintenance"},
  {key:"announcements_push",label:"Society announcements"},
  {key:"hub_replies_push",label:"Hub replies"},
];

export function NotificationPushPreferencesPanel(){
  const selectedFlat=useAppSelector((state)=>state.app.selectedFlatId);
  const selectedSociety=useAppSelector((state)=>state.app.selectedSocietyId);
  const bootstrap=useGetV1BootstrapQuery();
  const residences=getMobileResidences(bootstrap.data?.data);
  const memberships=getMobileMemberships(bootstrap.data?.data);
  const societyId=residences.find((r)=>r.flat_id===selectedFlat)?.society_id
    ?? selectedSociety ?? (residences.length===1?residences[0]?.society_id:undefined)
    ?? (memberships.length===1?memberships[0]?.society_id:undefined);
  const prefs=useGetNotificationPushPreferencesQuery(societyId ?? 0,{skip:!societyId});
  const [save,{isLoading}]=usePutNotificationPushPreferencesMutation();
  const [error,setError]=useState(false);
  if(!societyId||!prefs.data?.data)return null;
  const value=prefs.data.data;
  const update=async(key:keyof NotificationPushPreferences,enabled:boolean)=>{
    setError(false);
    try{await save({societyId,preferences:{...value,[key]:enabled}}).unwrap()}
    catch{setError(true)}
  };
  return <View style={styles.panel}>
    <Text style={styles.title}>Push preferences</Text>
    <Text style={styles.note}>Visitor approvals and important society alerts stay on. All updates remain in your inbox.</Text>
    {options.map(({key,label})=><View key={key} style={styles.row}>
      <Text style={styles.label}>{label}</Text>
      <Switch accessibilityLabel={`${label} push`} value={value[key]} disabled={isLoading} onValueChange={(enabled)=>{void update(key,enabled)}} />
    </View>)}
    {error?<Text style={styles.error}>Could not save preferences. Try again.</Text>:null}
  </View>;
}

const styles=StyleSheet.create({
  panel:{padding:16,marginBottom:12,borderRadius:12,backgroundColor:"#f3f6fa"},
  title:{fontSize:16,fontWeight:"700",marginBottom:4},
  note:{fontSize:12,color:"#526173",marginBottom:8},
  row:{minHeight:44,flexDirection:"row",alignItems:"center",justifyContent:"space-between"},
  label:{fontSize:14,flex:1},
  error:{fontSize:12,color:"#b42318",marginTop:8},
});
