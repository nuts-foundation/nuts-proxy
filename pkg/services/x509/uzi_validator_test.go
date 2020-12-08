/*
 * Nuts auth
 * Copyright (C) 2020. Nuts community
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package x509

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/nuts-foundation/nuts-auth/pkg/contract"
	"github.com/nuts-foundation/nuts-auth/pkg/services"
)

var uziSignedJwt = `eyJ4NWMiOlsiTUlJSGN6Q0NCVnVnQXdJQkFnSVVIUFU4cVZYS3FEZXByWUhDQ1dLQmkrdkp0Vll3RFFZSktvWklodmNOQVFFTEJRQXdhakVMTUFrR0ExVUVCaE1DVGt3eERUQUxCZ05WQkFvTUJFTkpRa2N4RnpBVkJnTlZCR0VNRGs1VVVrNU1MVFV3TURBd05UTTFNVE13TVFZRFZRUUREQ3BVUlZOVUlGVmFTUzF5WldkcGMzUmxjaUJOWldSbGQyVnlhMlZ5SUc5d0lHNWhZVzBnUTBFZ1J6TXdIaGNOTWpBd056RTNNVEl6TkRFNVdoY05Nak13TnpFM01USXpOREU1V2pDQmhURUxNQWtHQTFVRUJoTUNUa3d4SURBZUJnTlZCQW9NRjFURHFYTjBJRnB2Y21kcGJuTjBaV3hzYVc1bklEQXpNUll3RkFZRFZRUUVEQTEwWlhOMExUa3dNREUzT1RRek1Rd3dDZ1lEVlFRcURBTktZVzR4RWpBUUJnTlZCQVVUQ1Rrd01EQXlNVEl4T1RFYU1CZ0dBMVVFQXd3UlNtRnVJSFJsYzNRdE9UQXdNVGM1TkRNd2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUJEd0F3Z2dFS0FvSUJBUUNoVFloUEE3WDBTNWNWQnhHYzdHWi81RHZxSWVzaWowYUpadllMcVhrRmkzOU5EQjRLSDM4c3JIbHRGVWYyOVF3YlBSUm9KOEJJYXpFTnhkdTg4WUQvZXBKSGhmOUhpMkx1UGhoZmdSU3FjSnp4dDNPYStKME91YzdnZzBZaytnV01USkJ5R2ZSYlRQR3V5eVFFMnJOUFJteDRoOUNLSDZiNHVZam1ESDJWdXlhM3BtY0UrR2wxbmUvQnJjYnRsSmpCa2d6Vkw2cmVTYzdPUXhvbi9ZbmFRanhvakJpZ2xhT0hub2JESU9tczluQkZFQ29uUzVKNGZvb1VRVTg3anFMSGlHckJNL2xNdHlaOUVrblhGQ3U2U3VRb3ZDNlR1eUZ2c0JnT0MyNzNGZ0JaR2VybHkzbTFEVXczTlROUG15dlJEUXREWEJHTi9BVkVJLzR4VGdGL0FnTUJBQUdqZ2dMek1JSUM3ekJSQmdOVkhSRUVTakJJb0VZR0ExVUZCYUEvRmoweUxqRTJMalV5T0M0eExqRXdNRGN1T1RrdU1qRTRMVEV0T1RBd01ESXhNakU1TFU0dE9UQXdNREF6T0RJdE1EQXVNREF3TFRBd01EQXdNREF3TUF3R0ExVWRFd0VCL3dRQ01BQXdId1lEVlIwakJCZ3dGb0FVeWZBR0RwTGZOaThJZFRpODMrNUJlYkpkd0Y4d2dhc0dDQ3NHQVFVRkJ3RUJCSUdlTUlHYk1Hc0dDQ3NHQVFVRkJ6QUNobDlvZEhSd09pOHZkM2QzTG5WNmFTMXlaV2RwYzNSbGNpMTBaWE4wTG01c0wyTmhZMlZ5ZEhNdk1qQXhPVEExTURGZmRHVnpkRjkxZW1rdGNtVm5hWE4wWlhKZmJXVmtaWGRsY210bGNsOXZjRjl1WVdGdFgyTmhYMmN6TG1ObGNqQXNCZ2dyQmdFRkJRY3dBWVlnYUhSMGNEb3ZMMjlqYzNBdWRYcHBMWEpsWjJsemRHVnlMWFJsYzNRdWJtd3dnZ0VHQmdOVkhTQUVnZjR3Z2Zzd2dmZ0dDV0NFRUFHSGIyT0JWRENCNmpBL0JnZ3JCZ0VGQlFjQ0FSWXphSFIwY0hNNkx5OWhZMk5sY0hSaGRHbGxMbnB2Y21kamMzQXVibXd2WTNCekwzVjZhUzF5WldkcGMzUmxjaTVvZEcxc01JR21CZ2dyQmdFRkJRY0NBakNCbVF5QmxrTmxjblJwWm1sallXRjBJSFZwZEhOc2RXbDBaVzVrSUdkbFluSjFhV3RsYmlCMFpXNGdZbVZvYjJWMlpTQjJZVzRnWkdVZ1ZFVlRWQ0IyWVc0Z2FHVjBJRlZhU1MxeVpXZHBjM1JsY2k0Z1NHVjBJRlZhU1MxeVpXZHBjM1JsY2lCcGN5QnBiaUJuWldWdUlHZGxkbUZzSUdGaGJuTndjbUZyWld4cGFtc2dkbTl2Y2lCbGRtVnVkSFZsYkdVZ2MyTm9ZV1JsTGpBZkJnTlZIU1VFR0RBV0JnZ3JCZ0VGQlFjREJBWUtLd1lCQkFHQ053b0REREJqQmdOVkhSOEVYREJhTUZpZ1ZxQlVobEpvZEhSd09pOHZkM2QzTG5WNmFTMXlaV2RwYzNSbGNpMTBaWE4wTG01c0wyTmtjQzkwWlhOMFgzVjZhUzF5WldkcGMzUmxjbDl0WldSbGQyVnlhMlZ5WDI5d1gyNWhZVzFmWTJGZlp6TXVZM0pzTUIwR0ExVWREZ1FXQkJTWTBkclhRMEpINmhIdi9zejFTK3lyakVoU1F6QU9CZ05WSFE4QkFmOEVCQU1DQmtBd0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dJQkFGMDdXWmhoNkx5ZWdjMjJscDIwb0x5K2tnUlB3Ti9TL0lTdkxGVEY0RFBBSTY2RmtVSnNGUmFmbXVhMFpsL0JPZ2U1SXZwMHM5dEVqaHBaMTZYNGVZQm1qOE1VMHhBTjM0OC9PakFtSUZTR0l1d2kxU2RyendIUnF2VUxmMHNWcXZUOEpEVTZkMHEvaVBPRThEYU9OWXppbUlkZ1dFOXBOODhBb1ptT3VkSDQzSjk3WkRnMXYrWnU3NnMwdFI4WXpXSElUVDEvbmJRbDUzeU9mR3dER1RSdk42T1hkelBMVXpUbGhmdEdYZUZPRmNrb0Q4c2NRTGFaV1loQTVaVDRxLzlncE02WXU1TTMzWVJ0empGek4yTWVWaFpsUmV5NUY1NmVWcDV6MkM0U3NnM2FCemkyandnRzExY3pvMVBGdldod21zckNTTFpJUHdhWFduQ3hnYW5FZkxzeXVKcmpuVXYyUXdaeldCT1VoRjhSN2FtUk9xUHN6VGJwNE9yZWUyWmFyc04wYzNSLzdYdmJvcVdhb3NRa3Q1MFlxOHpCQ0Z4clFMZkZKN1pUcEhHWENEQmtzcVg4WWVrZ2RxdDhIMmdSS2p2OVNLY2RjejA0a2VJUEIyRU85K2ZQTHcwckZqRGVLdFFjYmRXTDlFSHRNOHAwcXBmTHNLcUdqbXdSdHhYbVRYUHNVS0FKQ1RKdWI4cnVRZVpsQlhZVC91YjNEMER1RzB2YUlNcjE3aDZydEdYR1hDWFV2VUxYMzBnczFyS3VUVkZkR0xFRUdid3JHbFVUZUdHRXFQbU4xdWFmNWpEdkR1UDE5R2RTV0VZMW4xTjYvV1paODhVS2ZnZHpxSVlKemt1RzV6bGZLUWdEREJvZXNyd3BCZXlkTXo0M0diZEZieS8zUm9MNSJdLCJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJtZXNzYWdlIjoiTkw6QmVoYW5kZWxhYXJMb2dpbjp2MSBPbmRlcmdldGVrZW5kZSBnZWVmdCB0b2VzdGVtbWluZyBhYW4gRGVtbyBFSFIgb20gbmFtZW5zIHZlcnBsZWVnaHVpcyBEZSBub290amVzIGVuIG9uZGVyZ2V0ZWtlbmRlIGhldCBOdXRzIG5ldHdlcmsgdGUgYmV2cmFnZW4uIERlemUgdG9lc3RlbW1pbmcgaXMgZ2VsZGlnIHZhbiBkaW5zZGFnLCAxIG9rdG9iZXIgMjAxOSAxMzozMDo0MiB0b3QgZGluc2RhZywgMSBva3RvYmVyIDIwMTkgMTQ6MzA6NDIuIiwiaWF0IjoxNjA0MzE3ODg5fQ.FMekUy0UoOwhbEciJ9Q1TESh7fE-MQuUEZI5M65RuwtTlPlqN2P1KGFel8FDh42k2R79S8RB4x1XF0UkZtu8YOkNqFuX2h5Ow3xhaAquHR3iqzJy8wBKo0ZnctPDSJGfn0k-UzF9MS6665JuDAnvE5ETop1ASou2lPC6885Rh8QRxBDSKz48pHsLh2oQrn7Qs5BfhHMgkDrwnPrN1tIhyKPNvbhFvy7nYbrdKg6O3W8xK9jHyES7ts_ahkI3GYH9nOa2VhX3lySLzsY3qH5NPDNCj3IE1St6Ab4rm7RfCQ8tWVRf0qQG1X0bALgCNMY8ALUrIoUUn4zxpAGCNRBmig`

func TestUziValidator(t *testing.T) {
	t.Run("ok - it loads the production certificates", func(t *testing.T) {
		_, err := NewUziValidator(UziProduction, &contract.StandardContractTemplates, nil)
		if !assert.NoError(t, err) {
			return
		}
	})

	t.Run("ok - acceptation environment", func(t *testing.T) {
		crls, err := NewMockCrlService([]string{
			"http://www.uzi-register-test.nl/cdp/test_uzi-register_medewerker_op_naam_ca_g3.crl",
			"http://www.uzi-register-test.nl/cdp/test_zorg_csp_level_2_persoon_ca_g3.crl",
			"http://www.uzi-register-test.nl/cdp/test_zorg_csp_root_ca_g3.crl"})

		if !assert.NoError(t, err) {
			return
		}
		uziValidator, err := NewUziValidator(UziAcceptation, &contract.StandardContractTemplates, crls)
		if !assert.NoError(t, err) {
			return
		}

		signedToken, err := uziValidator.Parse(uziSignedJwt)

		if !assert.NoError(t, err) {
			return
		}

		if !assert.Implements(t, (*services.SignedToken)(nil), signedToken) {
			return
		}

		expected := map[string]string{
			"agbCode":  "00000000",
			"cardType": "N",
			"oidCa":    "2.16.528.1.1007.99.218", // CIBG.Uzi test identifiers
			"orgID":    "90000382",
			"rollCode": "00.000",
			"uziNr":    "900021219",
			"version":  "1",
		}
		attrs, err := signedToken.SignerAttributes()

		assert.NoError(t, err)
		assert.Equal(t, expected, attrs)
		assert.Equal(t, contract.Type("BehandelaarLogin"), signedToken.Contract().Template.Type)
		assert.Equal(t, contract.Language("NL"), signedToken.Contract().Template.Language)
		assert.Equal(t, contract.Version("v1"), signedToken.Contract().Template.Version)

		// Replace the time func with one that returns a time the crl is valid
		oldNowFunc := nowFunc
		defer func() {
			nowFunc = oldNowFunc
		}()
		nowFunc = func() time.Time { return time.Date(2020, 10, 29, 0, 0, 0, 0, time.UTC) }

		err = uziValidator.Verify(signedToken)
		assert.NoError(t, err)
	})
}
