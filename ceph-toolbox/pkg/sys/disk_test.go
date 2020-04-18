package sys

import "testing"

const diskInfo = `文件系统           1K-块      已用      可用 已用% 挂载点
udev             8069088         0   8069088    0% /dev
tmpfs            1618656      2376   290141484    1% /run
/dev/nvme0n1p2 490691512 175554548 1616280   38% /
tmpfs            8093268    280096   7813172    4% /dev/shm
tmpfs               5120         4      5116    1% /run/lock
tmpfs            8093268         0   8093268    0% /sys/fs/cgroup
/dev/nvme0n1p1    523248      6152    517096    2% /boot/efi
/dev/loop0         10624     10624         0  100% /snap/kompose/1
/dev/loop1          1024      1024         0  100% /snap/gnome-logs/57
/dev/loop2         43904     43904         0  100% /snap/gtk-common-themes/1313
/dev/loop3        144128    144128         0  100% /snap/gnome-3-26-1604/92
/dev/loop4        153600    153600         0  100% /snap/gnome-3-28-1804/67
/dev/loop7          1024      1024         0  100% /snap/gnome-logs/61
/dev/loop5          4224      4224         0  100% /snap/gnome-calculator/406
/dev/loop6         55808     55808         0  100% /snap/core18/1074
/dev/loop8        153600    153600         0  100% /snap/gnome-3-28-1804/71
/dev/loop9         91264     91264         0  100% /snap/core/7713
/dev/loop10         3840      3840         0  100% /snap/gnome-system-monitor/95
/dev/loop12       135424    135424         0  100% /snap/postman/81
/dev/loop11        21504     21504         0  100% /snap/gnome-logs/25
/dev/loop13        15104     15104         0  100% /snap/gnome-characters/296
/dev/loop15        90880     90880         0  100% /snap/core/7396
/dev/loop14       153216    153216         0  100% /snap/slack/16
/dev/loop16       144128    144128         0  100% /snap/gnome-3-26-1604/90
/dev/loop17        15104     15104         0  100% /snap/gnome-characters/317
/dev/loop18        10240     10240         0  100% /snap/helm/124
/dev/loop19         4352      4352         0  100% /snap/gnome-calculator/501
/dev/loop20         8704      8704         0  100% /snap/canonical-livepatch/81
/dev/loop21       150144    150144         0  100% /snap/slack/17
/dev/loop22         3840      3840         0  100% /snap/gnome-system-monitor/100
/dev/loop23        55808     55808         0  100% /snap/core18/1098
/dev/loop24       148480    148480         0  100% /snap/notepadqq/855
/dev/loop25        36224     36224         0  100% /snap/gtk-common-themes/1198
tmpfs            1618652        20   1618632    1% /run/user/120
tmpfs            1618652        36   1618616    1% /run/user/1000`

func TestLookUpValidDisk(t *testing.T) {
	exp := "/run"
	act, _ := LookUpValidDisk(diskInfo)
	if exp != act {
		t.Errorf("exp: %s act: %s \n", exp, act)
	}
}
