package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"gorm.io/gorm/clause"
)

const WitnessContent = `Content-Type: multipart/mixed; boundary="===============4504313173045938348=="
MIME-Version: 1.0
From: Gitolite <devnull@kernel.org>
Subject: post-receive: pub/scm/linux/kernel/git/gregkh/driver-core
Date: Sun, 02 May 2021 15:38:12 -0000
Message-Id: <161996989292.9223.15453726891360891741@gitolite.kernel.org>

--===============4504313173045938348==
Content-Type: text/plain; charset="us-ascii"
MIME-Version: 1.0
Content-Transfer-Encoding: 7bit

---
service: git-receive-pack
repo: pub/scm/linux/kernel/git/gregkh/driver-core
user: gregkh
git_push_cert_status: G
changes:
  - ref: refs/heads/umn.edu-reverts
    old: 5094fcc2277a53380510c0f313ec8daf8bab3a5a
    new: 2c913b8ec4971f82c26a5828d2b887849258df93
    log: revlist-5094fcc2277a-2c913b8ec497.txt

--===============4504313173045938348==
Content-Type: text/plain; charset="us-ascii"
MIME-Version: 1.0
Content-Transfer-Encoding: 7bit
Content-Disposition: attachment; filename=git-push-certificate.txt

certificate version 0.1
pusher Greg Kroah-Hartman <gregkh@linuxfoundation.org> 1619969884 +0200
pushee gitolite.kernel.org:/pub/scm/linux/kernel/git/gregkh/driver-core.git
nonce 1619969883-e339dfdee5a4cd1adcc3dccedf5361dd9a7e2384

5094fcc2277a53380510c0f313ec8daf8bab3a5a 2c913b8ec4971f82c26a5828d2b887849258df93 refs/heads/umn.edu-reverts
-----BEGIN PGP SIGNATURE-----

iQJPBAABCAA5FiEEZH8oZUiU471FcZm+ONu9yGCSaT4FAmCOx1wbHGdyZWdraEBs
aW51eGZvdW5kYXRpb24ub3JnAAoJEDjbvchgkmk+QN8P/3D+gl7s5/YpW9tSxipi
m2eRS5zSdHJ+t7j0F02ulUwUdI8XnG8CZIRevCiCiMy6tJDK9J96B6dAv/q4u88v
6VInFzl0lW75zGyS7BTMz/rSDA/M5GgmOadrSZJaBbZyO27SYyO4wae8s+7cgGz6
H9lj0+zXS24y+uZeQfjkVptB1oPGZ0pkJ0M+S6ir89y5W5H7caIKxvqFZ+0xsncK
vEC5d2PUaXYCtQZDykpgikU7L6RtAFiV54i+c2RPyyr++g0jRSYaxfHBazKBqoU5
kFlo5RMP17epzgNEPI4ZVtgj8IukItnIvVOW3h5wxU2JydKvj37X2/IiXJbnATw2
6Rvhr3xsThujcX0WUApvGVUWInnngeEgbe8EEMEIGkWO86l3EuixvA2ZDSvN11pM
+3jlq9EMLtlcE9hN00n2cWkCVcbaDwGgYcJh2nUiHXgIn6fEsIGRvHoDSYDQNh+i
+4ksj2skeof/9gs2WWAaWDpUzujoAEa91E0bZQCxZacVIbKREwbUERCPOiE234or
CwsbYVVl0BgPoYmCYA88v/q3bo78JP7R3Les3nId1ZONHPo1Qo4qy9uBw8PiX9Mf
Fxt9fiTV0hheEB5SFxYT3b0gIIrLb2IYzDlFOnLhzR93rBc3AF8Ut+0WB3sswefF
CJ6+INkL7qtP+VvXJoWs0nAG
=/YUO
-----END PGP SIGNATURE-----

--===============4504313173045938348==
Content-Type: text/plain; charset="us-ascii"
MIME-Version: 1.0
Content-Transfer-Encoding: 7bit
Content-Disposition: attachment; filename=revlist-5094fcc2277a-2c913b8ec497.txt

c1a4a3326d55192442ff9cd4f4a08c6fdaa44e0d Revert "ALSA: usx2y: Fix potential NULL pointer dereference"
b9bc317ba346288bc5ec21261ad925fa30f83b72 Revert "ALSA: sb: fix a missing check of snd_ctl_add"
1627e5b20f7fbc7d11a7dd517642c266204a8048 Revert "hwmon: (lm80) fix a missing check of bus read in lm80 probe"
494b23fb5e7f7ed19b16da2ae4b5ed1a921fb00e Revert "video: imsttfb: fix potential NULL pointer dereferences"
1040b469d861bff498faeb74fae2e9a6bdf51e30 Revert "brcmfmac: add a check for the status of usb_register"
75a89c9de383085f14f9305064005b20419440fe Revert "ASoC: cs43130: fix a NULL pointer dereference"
c31095c22f65c32d8d96c4d76fa62bce577b9e9e Revert "ASoC: rt5645: fix a NULL pointer dereference"
4be1d7e9ebad3a6cf833f8f849dd2e929bd93325 Revert "net: liquidio: fix a NULL pointer dereference"
1ed92dfdd0b20d4571e82ef2df38f23c41768df1 !!!! Canary - fixes
820e5480ad918bfa1aed90dd6449e548807bd68e Revert "crypto: cavium/nitrox - add an error message to explain the failure of pci_request_mem_regions"
a02129269ebbd9468fdfcb000611b0cde8f41732 Revert "media: rcar_drif: fix a memory disclosure"
507046cf93099317b4bab958c0dc9beb59663806 Revert "rtlwifi: fix a potential NULL pointer dereference"
ee5867e1b8757bf422e4e63cf5a36912720f89ec net: rtlwifi: properly check for alloc_workqueue() failure
db23849c217c2531a333d21c2b89432acf24989d Revert "net: fujitsu: fix a potential NULL pointer dereference"
88980b41038e59809bd33fc8b6f9e61488cfce15 net: fujitsu: fix potential null-ptr-deref
cd4adb94f039f4499ea7f36e9fdcc044e4760e27 Revert "leds: lp5523: fix a missing check of return value of lp55xx_read"
7af64b93972184a504667b74c90dbda2a568fc1c leds: lp5523: check return value of lp5xx_read and jump to cleanup code
48e750ef9fbaba40cb8d2e7ffa57f0d1e7d65caf Revert "serial: max310x: pass return value of spi_register_driver"
7e8bb3667106ba4723ffb6e9ad2395de0ddf2b3a serial: max310x: unregister uart driver in case of failure and abort
e5b622e57b2b08f20828ead059db36338c6a54b0 Revert "net/smc: fix a NULL pointer dereference"
19e2cc0daf4191fa2261da3ee986abe3bdf2aa38 net/smc: properly handle workqueue allocation failure
99f209660e9b14dc50858641c876dcbd78d5fab2 Revert "char: hpet: fix a missing check of ioremap"
25293d9adcffde9fdcbd28538834a8c983999b6c char: hpet: add checks after calling ioremap
03a566bfc203219cae3cd03def7392beab7d8b1a Revert "scsi: ufs: fix a missing check of devm_reset_control_get"
b2010469a58d8c208b3bf54573e2796988ee4a8c scsi: ufs: handle cleanup correctly on devm_reset_control_get error
8c88aeb95ca67f9d9fd2af8e4c2a38bbb2fb0e47 Revert "net: caif: replace BUG_ON with recovery code"
db7a9bbb7e59f3954e82e94b45e04ec7135c2f4a net: caif: remove BUG_ON(dev == NULL) in caif_xmit
0aa4e38e9011bd849b10f9e8e0d6b737e26bbb03 Revert "net: stmicro: fix a missing check of clk_prepare"
4663d5455e8e44f2bcfd398049dd9ce4675e5c40 net: stmicro: handle clk_prepare() failure during init
0b564074ff0b6d1ce5914f7b6cb971679bcf65ec Revert "gdrom: fix a memory leak bug"
da1af8571a14877af99bb6383adfe0a43aafb0dc cdrom: gdrom: deallocate struct gdrom_unit fields in remove_gdrom
5be33acaabda108c1ebef623fea8b612d24365f8 Revert "ALSA: gus: add a check of the status of snd_ctl_add"
6c6bf12179fd4ec693c384ee17bc9d1e1d969d42 ALSA: gus: properly handle snd_ctl_add() error
41d6e07c3a177f130d0fb5d8764976193dffbd4b Revert "video: hgafb: fix potential NULL pointer dereference"
baffe9bf720574bf5bfc4d66886254c00527962a video: hgafb: fix potential NULL pointer dereference
f7703e2c617af0a332a1f3238677ea3f87974819 Revert "isdn: mISDNinfineon: fix potential NULL pointer dereference"
bc3a58d11e3f770af70180d651e4fd2f17570e44 isdn: mISDNinfineon: check/cleanup ioremap failure correctly in setup_io
a4857e145869d7ce1b2464e8092f744d922cc903 Revert "ath6kl: return error code in ath6kl_wmi_set_roam_lrssi_cmd()"
970aa11fb9df0333c02b5544f02c44bd590d9c69 ath6kl: return error code in ath6kl_wmi_set_roam_lrssi_cmd()
63a8a54e9901de40a433f2e30d4437eea7356882 Revert "rapidio: fix a NULL pointer dereference when create_workqueue() fails"
050e0ab3d55d83dfb10e9a09d2791d08dd561044 rapidio: handle create_workqueue() failure
f233439db08ee44e05364bf219f13b164ec16c76 Revert "niu: fix missing checks of niu_pci_eeprom_read"
e5135579d19b1431b67dc4bf22e72b43610c5fdf ethernet: sun: niu: fix missing checks of niu_pci_eeprom_read()
24458cbb508d53fd7dd6a6daf62e4695a54bc0bf Revert "isdn: mISDN: Fix potential NULL pointer dereference of kzalloc"
81fdd3be17d50c724968e3a8330185bcb2fcafd2 isdn: mISDN: correctly handle ph_info allocation failure in hfcsusb_ph_info
0349db04824dbd9a8b6f971d8b2b919a7d19ab13 Revert "ecryptfs: replace BUG_ON with error handling code"
cfa8b73d965463707b8b650c442852d216bf23c0 fs: ecryptfs: remove BUG_ON from crypt_scatterlist
59f30493bf8b202de048c16ac082b5a2dea112b0 Revert "dmaengine: qcom_hidma: Check for driver register failure"
8260fe0b6f989af9c940cf8a06e156f953355d3f dmaengine: qcom_hidma: comment platform_driver_register call
0d56e04bd4988cc4a37f25e02376ae32b29a3fbf Revert "qlcnic: Avoid potential NULL pointer dereference"
bdc27019fbe872da1448555f8545afa01d92354c qlcnic: Add null check after calling netdev_alloc_skb
7a2673489ba49ad0766f24840023e8b2effe1d1f Revert "libertas: add checks for the return value of sysfs_create_group"
f1cadd1ab834b9daf69be2dac39af393f1cd606c libertas: register sysfs groups properly
8ded1ce22d596bffd9346cf865f5c730b0760ff7 Revert "ALSA: sb8: add a check for request_region"
2c913b8ec4971f82c26a5828d2b887849258df93 ALSA: sb8: Add a comment note regarding an unused pointer

--===============4504313173045938348==--
`

func TestRevisions(t *testing.T) {
	r := strings.NewReader(WitnessContent)
	files, _ := ParseMail(r)
	w, _ := WorkRevisions(files)
	fmt.Println(len(w))
}

func TestWorkDB(t *testing.T) {
	InitDB("test.db")
	r := strings.NewReader(WitnessContent)
	files, _ := ParseMail(r)
	tlogEntry := WorkTLog("2c913b8ec4971f82c26a5828d2b887849258df93", files)
	AddCommit(tlogEntry)
	var entry Revision
	db.Preload(clause.Associations).Find(&entry)
	e, _ := json.MarshalIndent(entry, " ", " ")
	fmt.Println(string(e))
}
